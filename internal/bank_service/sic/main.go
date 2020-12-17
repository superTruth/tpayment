package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/internal/bank_service/bank_common"
	"tpayment/internal/bank_service/bank_common/api"
	"tpayment/internal/bank_service/sic/internal/common"
	"tpayment/internal/bank_service/sic/internal/logon"
	"tpayment/internal/bank_service/sic/internal/sale"
	"tpayment/pkg/algorithmutils"
	"tpayment/pkg/emv/tlv"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils/convert_utils"
	"tpayment/pkg/utils/format_utils"

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

	// 查看是否有TMK
	tmk := common.FindKey(reqBean, bank_common.TMK)
	if tmk == nil { // 需要下载主秘钥
		logger.Info("download PUK")
		// step 1： 下载公钥
		downloadPuk := logon.DownloadPuk{
			Req:    reqBean,
			Config: config,
		}
		_, errorCode := downloadPuk.Handle()
		if errorCode != conf.Success {
			logger.Warn("downloadPuk.Handle fail ->", errorCode.String())
			baseResp.ErrorCode = string(errorCode)
			baseResp.ErrorDes = errorCode.String()
			return baseResp, nil
		}

		// step 2: 下载主秘钥
		logger.Info("download TMK")
		downloadTmk := logon.DownloadTMK{
			Req:         reqBean,
			Config:      config,
			PukModulus:  downloadPuk.PukModulus,
			PukExponent: downloadPuk.PukExponent,
		}
		resp, errorCode := downloadTmk.Handle()
		if errorCode != conf.Success {
			logger.Warn("downloadTmk.Handle fail ->", errorCode.String())
			baseResp.ErrorCode = string(errorCode)
			baseResp.ErrorDes = errorCode.String()
			return baseResp, nil
		}
		logger.Info("continue transaction")
		baseResp.ErrorCode = string(conf.NeeContinue)
		baseResp.ErrorDes = conf.NeeContinue.String()
		baseResp.RespBody = formatResponse(resp)
		return baseResp, nil
	}

	// 查看是否包含WK
	tdk := common.FindKey(reqBean, bank_common.TDK)
	if tdk == nil {
		logger.Info("download WK")
		// step 1： 下载WK
		downloadPuk := logon.DownloadWK{
			Req:    reqBean,
			Config: config,
		}
		resp, errorCode := downloadPuk.Handle()
		if errorCode != conf.Success {
			logger.Warn("downloadPuk.Handle fail ->", errorCode.String())
			baseResp.ErrorCode = string(errorCode)
			baseResp.ErrorDes = errorCode.String()
			return baseResp, nil
		}
		logger.Info("continue transaction")
		baseResp.ErrorCode = string(conf.NeeContinue)
		baseResp.ErrorDes = conf.NeeContinue.String()
		baseResp.RespBody = formatResponse(resp)
		return baseResp, nil
	}

	// 加密敏感数据
	err = enImportantData(reqBean)
	if err != nil {
		logger.Error("encrypt important data fail->", err.Error())
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

func formatResponse(resp *acquirer_impl.SaleResponse) string {
	retBytes, _ := json.Marshal(resp)
	return string(retBytes)
}

func enImportantData(req *acquirer_impl.SaleRequest) error {
	if len(req.TxqReq.CreditCardBean.PIN) != 0 { // 加密PIN
		tpk := common.FindKey(req, bank_common.TPK)
		if tpk == nil {
			return errors.New("No Tpk")
		}
		if len(req.TxqReq.CreditCardBean.CardNumber) == 0 {
			return errors.New("No Card No")
		}

		tpkEn, err := algorithmutils.EncryptPIN(req.TxqReq.CreditCardBean.PIN,
			req.TxqReq.CreditCardBean.CardNumber, convert_utils.HexString2Bytes(tpk.Value))
		if err != nil {
			return err
		}
		req.TxqReq.CreditCardBean.PIN = convert_utils.Bytes2HexString(tpkEn)
	}

	if len(req.TxqReq.CreditCardBean.CardTrack2) != 0 { // 加密TK2
		tdk := common.FindKey(req, bank_common.TDK)

		// 去除tk2结尾的F
		tk2Bytes := []byte(strings.ToUpper(req.TxqReq.CreditCardBean.CardTrack2))
		if tk2Bytes[len(tk2Bytes)-1] == 'F' {
			req.TxqReq.CreditCardBean.CardTrack2 =
				req.TxqReq.CreditCardBean.CardTrack2[:len(req.TxqReq.CreditCardBean.CardTrack2)-1]
		}

		// 添加长度
		formatTk2 := fmt.Sprintf("%2d", len(req.TxqReq.CreditCardBean.CardTrack2)) + (req.TxqReq.CreditCardBean.CardTrack2)

		// 格式化成8的整数倍
		formatTk2 = format_utils.AppendString(formatTk2, (len(formatTk2)+15)/16*16, false, '0')

		formatTk2Bytes := convert_utils.HexString2Bytes(formatTk2)

		enFormatedTk2Bytes, err := algorithmutils.EncryptDesECB(formatTk2Bytes, convert_utils.HexString2Bytes(tdk.Value))

		if err != nil {
			return err
		}
		req.TxqReq.CreditCardBean.CardTrack2 = convert_utils.Bytes2HexString(enFormatedTk2Bytes)
	}

	return nil
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
