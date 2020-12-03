package standard

import (
	"encoding/json"
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/internal/bank_service/bank_common"
	"tpayment/pkg/grpc_pool"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func (api *API) Sale(ctx *gin.Context, req *acquirer_impl.SaleRequest) (*acquirer_impl.SaleResponse, conf.ResultCode) {
	//logger := tlog.GetLogger(ctx)
	resp := new(acquirer_impl.SaleResponse)
	resp.TxnResp = new(api_define.TxnResp)

	var (
		errorCode conf.ResultCode
	)

	// 参数检查
	errorCode = saleValidate(ctx, req)
	if errorCode != conf.Success {
		return nil, errorCode
	}

	//// 在线交易
	//bankResp, errorCode := saleOnline(ctx, req)
	//if errorCode != conf.Success {
	//	logger.Warn("online fail->", errorCode.String())
	//	return resp, errorCode
	//}
	//
	//resp.TxnResp.CreditCardBean = &api_define.CreditCardBean{
	//	AuthCode:     bankResp.TxnResp.CreditCardBean.AuthCode,
	//	ResponseCode: bankResp.TxnResp.CreditCardBean.ResponseCode,
	//	IccResponse:  bankResp.TxnResp.CreditCardBean.IccResponse,
	//}
	//resp.TxnResp.AcquirerRRN = bankResp.TxnResp.AcquirerRRN
	//resp.AcquirerReconID = bankResp.AcquirerReconID
	//
	//resp.TxnResp.TransactionState = record.Success
	//resp.TxnResp.AdditionData = bankResp.TxnResp.AdditionData

	resp.TxnResp.CreditCardBean = &api_define.CreditCardBean{
		AuthCode:     "1234",
		ResponseCode: "1234",
	}

	return resp, conf.Success
}

func saleValidate(ctx *gin.Context, req *acquirer_impl.SaleRequest) conf.ResultCode {
	logger := tlog.GetLogger(ctx)

	if req.TxqReq.PaymentProcessRule.MerchantAccount.Terminal == nil {
		logger.Warn("not match tid->", req.TxqReq.DeviceID)
		return conf.NotMatchTID
	}

	return conf.Success
}

func saleOnline(ctx *gin.Context, req *acquirer_impl.SaleRequest) (*acquirer_impl.SaleResponse, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)
	var (
		err error
	)
	// 从配置上面获取到微服务连接方式
	acqAddition := req.TxqReq.PaymentProcessRule.MerchantAccount.Acquirer.Addition
	acquirerConfig := new(AcquirerConfigDefault)
	err = json.Unmarshal([]byte(acqAddition), acquirerConfig)
	if err != nil {
		logger.Error("can't parse acquirer addition->", acqAddition)
		return nil, conf.ConfigError
	}
	if acquirerConfig.GRPCConnectInfo == "" {
		logger.Error("can't get service GRPCConnectInfo from acquirer addition->", acqAddition)
		return nil, conf.ConfigError
	}

	// 获取连接对象
	conn, err := grpc_pool.GetConn(acquirerConfig.GRPCConnectInfo)
	if err != nil {
		logger.Error("can't reach bank services->", acquirerConfig.GRPCConnectInfo)
		return nil, conf.CantReachAcquirer
	}
	defer grpc_pool.PutConn(acquirerConfig.GRPCConnectInfo, conn)

	// 序列化请求数据
	reqByte, _ := json.Marshal(req)
	bankReq := &bank_common.BaseRequest{
		ReqBody: string(reqByte),
	}

	c := bank_common.NewTxnClient(conn)
	bankResp, err := c.BaseTxn(ctx, bankReq)
	if err != nil {
		logger.Error("c.BaseTxn grpc fail->", err.Error())
		return nil, conf.Reversal
	}

	// 解析返回数据
	if bankResp.ErrorCode != string(conf.Success) {
		return nil, conf.ResultCode(bankResp.ErrorCode)
	}

	bankRespBody := new(acquirer_impl.SaleResponse)
	err = json.Unmarshal([]byte(bankResp.RespBody), bankRespBody)
	if err != nil {
		logger.Error("can't parse bank service response")
		return nil, conf.Reversal
	}

	//

	return bankRespBody, conf.Success
}
