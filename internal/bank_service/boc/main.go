package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/internal/bank_service/bank_common/api"
	"tpayment/internal/bank_service/boc/internal/common"
	"tpayment/internal/bank_service/boc/internal/sale"
	"tpayment/pkg/emv/tlv"
	"tpayment/pkg/tlog"

	"github.com/google/uuid"

	"google.golang.org/grpc"
)

type Server struct {
}

func main() {
	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		fmt.Println("net.Listen fail->", err.Error())
		return
	}

	s := grpc.NewServer()

	api.RegisterTxnServer(s, &Server{})
	err = s.Serve(lis)
	if err != nil {
		panic(err.Error())
	}
}

func (s *Server) EmptyCall(context.Context, *api.EmptyMessage) (*api.EmptyMessage, error) {
	return &api.EmptyMessage{}, nil
}

func (s *Server) BaseTxn(ctx context.Context, in *api.BaseRequest) (*api.BaseReply, error) {
	// 日志初始化
	requestId := uuid.New().String()
	logger := tlog.NewLog(requestId)
	tlog.SetGoroutineLogger(logger)
	defer tlog.FreeGoroutineLogger()

	logger.Info("BaseTxn->", in.ReqBody)

	reqBean := new(acquirer_impl.SaleRequest)
	baseResp := new(api.BaseReply)
	baseResp.ErrorCode = string(conf.Success)

	// 解包数据
	err := json.Unmarshal([]byte(in.ReqBody), reqBean)
	if err != nil {
		logger.Error("json.Unmarshal error->", err.Error())
		baseResp.ErrorCode = string(conf.ParameterError)
		baseResp.ErrorDes = conf.ParameterError.String()
		return baseResp, nil
	}

	// 解析config数据
	config, err := common.ParseConfig(reqBean)
	if err != nil {
		logger.Error("ParseConfig error->", err.Error())
		baseResp.ErrorCode = string(conf.ParameterError)
		baseResp.ErrorDes = conf.ParameterError.String()
		return baseResp, nil
	}

	// 过滤处理ICC Data
	filterIccData(reqBean)

	// 分类交易
	var (
		resp      *acquirer_impl.SaleResponse
		errorCode conf.ResultCode
	)
	logger.Info("start transaction")
	switch reqBean.TxqReq.TxnType {
	case conf.Sale:
		saleAction := &sale.Sale{
			Req:    reqBean,
			Config: config,
		}
		resp, errorCode = saleAction.Handle()
	default:
		baseResp.ErrorCode = string(conf.NotSupport)
		baseResp.ErrorDes = conf.NotSupport.String()
		return baseResp, nil
	}

	respBody, _ := json.Marshal(resp)
	baseResp.RespBody = string(respBody)
	baseResp.ErrorCode = string(errorCode)

	logger.Info("transaction ret code->", errorCode, ", body->", baseResp.RespBody)

	return baseResp, nil
}

func filterIccData(req *acquirer_impl.SaleRequest) {
	if req.TxqReq.CreditCardBean.IccRequest == "" {
		return
	}

	// 删除不必要的字段
	mapData, err := tlv.Parse2Map(req.TxqReq.CreditCardBean.IccRequest, false)
	if err != nil {
		return
	}
	if _, ok := mapData["5A"]; !ok {
		return
	}
	delete(mapData, "5A")
	req.TxqReq.CreditCardBean.IccRequest = tlv.Format(mapData)
}
