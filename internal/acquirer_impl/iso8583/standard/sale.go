package standard

import (
	"strconv"
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/models"
	"tpayment/models/payment/acquirer"
	"tpayment/models/payment/record"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func (wlb *API) Sale(ctx *gin.Context, req *acquirer_impl.SaleRequest) (*acquirer_impl.SaleResponse, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)
	resp := new(acquirer_impl.SaleResponse)
	resp.TxnResp = new(api_define.TxnResp)

	var (
		err       error
		errorCode conf.ResultCode
	)

	// 参数检查
	errorCode = saleValidate(ctx, req)
	if errorCode != conf.Success {
		return nil, errorCode
	}

	// 获取账号信息
	account := &acquirer.Account{
		BaseModel: models.BaseModel{
			Db:  models.DB(),
			Ctx: ctx,
		},
	}
	account, err = account.GetOrCreate(
		strconv.Itoa(int(req.TxqReq.PaymentProcessRule.MerchantAccountID)),
		GetAccountTag(req))
	if err != nil {
		logger.Error("account.GetOrCreate error->", err.Error())
		return nil, conf.DBError
	}
	account.Db = models.DB()
	account.Ctx = ctx

	// 拼接发送数据
	req.TxqReq.CreditCardBean.TraceNum = account.TraceNum
	req.TxqReq.CreditCardBean.BatchNum = account.BatchNum

	// 流水号增加
	err = account.IncTraceNum()
	if err != nil {
		logger.Error("IncTraceNum ", account.Tag, "fail->", err.Error())
		return nil, conf.DBError
	}

	resp.TxnResp.CreditCardBean = &api_define.CreditCardBean{
		TraceNum:     req.TxqReq.CreditCardBean.TraceNum,
		BatchNum:     req.TxqReq.CreditCardBean.BatchNum,
		AuthCode:     "1234",
		ResponseCode: "00",
	}

	// 开始交易
	resp.TxnResp.TransactionState = record.Success

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
