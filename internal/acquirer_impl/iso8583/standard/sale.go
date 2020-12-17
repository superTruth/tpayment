package standard

import (
	"encoding/json"
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/internal/bank_service/bank_common/api"
	"tpayment/models"
	"tpayment/models/payment/acquirer"
	"tpayment/pkg/grpc_pool"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func (api *API) Sale(ctx *gin.Context, req *acquirer_impl.SaleRequest) (*acquirer_impl.SaleResponse, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)

	var err error

	resp := new(acquirer_impl.SaleResponse)
	resp.TxnResp = new(api_define.TxnResp)

	// 参数检查
	var errorCode = saleValidate(ctx, req)
	if errorCode != conf.Success {
		return nil, errorCode
	}

	// 查找所有Key
	keyBean := acquirer.Key{
		BaseModel: models.BaseModel{
			Db: models.DB(),
		},
	}
	keyTag := generateKeyTag(req)
	req.Keys, err = keyBean.Get(keyTag)
	if err != nil {
		logger.Error("keyBean.Get fail->", err.Error())
		return nil, conf.DBError
	}

	// 在线交易
	bankResp, errorCode := saleOnline(ctx, req)
	logger.Info("services saleOnline result->", errorCode)
	// 需要注入key
	if bankResp != nil && len(bankResp.Keys) != 0 {
		err := updateKey(ctx, req, bankResp)
		if err != nil {
			logger.Error("updateKey fail->", err.Error())
			return nil, conf.DBError
		}
	}

	if errorCode != conf.Success {
		return resp, errorCode
	}

	if bankResp == nil || bankResp.TxnResp == nil {
		logger.Error("TxnResp is null")
		return resp, conf.UnknownError
	}

	if bankResp.TxnResp.CreditCardBean != nil {
		resp.TxnResp.CreditCardBean = &api_define.CreditCardBean{
			AuthCode:     bankResp.TxnResp.CreditCardBean.AuthCode,
			ResponseCode: bankResp.TxnResp.CreditCardBean.ResponseCode,
			IccResponse:  bankResp.TxnResp.CreditCardBean.IccResponse,
		}
	}
	resp.TxnResp.AcquirerRRN = bankResp.TxnResp.AcquirerRRN
	resp.AcquirerReconID = bankResp.AcquirerReconID

	resp.TxnResp.TransactionState = bankResp.TxnResp.TransactionState
	resp.TxnResp.ErrorCode = bankResp.TxnResp.ErrorCode
	resp.TxnResp.ErrorDesc = bankResp.TxnResp.ErrorDesc
	resp.TxnResp.AdditionData = bankResp.TxnResp.AdditionData

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
	reqByte, err := json.Marshal(req)
	if err != nil {
		logger.Error("json.Marshal fail->", err.Error())
		return nil, conf.ParameterError
	}
	bankReq := &api.BaseRequest{
		ReqBody: string(reqByte),
	}

	c := api.NewTxnClient(conn)

	// 测试连接
	_, err = c.EmptyCall(ctx, &api.EmptyMessage{})
	if err != nil {
		logger.Error("can't reach bank services->", acquirerConfig.GRPCConnectInfo)
		return nil, conf.CantReachAcquirer
	}

	logger.Info("bankReq body->", bankReq.ReqBody)
	bankResp, err := c.BaseTxn(ctx, bankReq)
	if err != nil {
		logger.Error("c.BaseTxn grpc fail->", err.Error())
		return nil, conf.Reversal
	}
	logger.Info("bankResp code->", bankResp.ErrorCode, ", body->", bankResp.RespBody)

	// 解析返回数据
	//if bankResp.ErrorCode != string(conf.Success) {
	//	return nil, conf.ResultCode(bankResp.ErrorCode)
	//}

	var bankRespBody *acquirer_impl.SaleResponse
	if len(bankResp.RespBody) != 0 {
		bankRespBody = new(acquirer_impl.SaleResponse)
		err = json.Unmarshal([]byte(bankResp.RespBody), bankRespBody)
		if err != nil {
			logger.Error("can't parse bank service response")
			return nil, conf.Reversal
		}
	}

	return bankRespBody, conf.ResultCode(bankResp.ErrorCode)
}
