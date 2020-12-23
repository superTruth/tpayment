package pay_manage

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/models/payment/acquirer"
	"tpayment/models/payment/merchantaccount"
	"tpayment/models/payment/paymentprocessrule"
	"tpayment/models/payment/record"
	"tpayment/modules"
	"tpayment/modules/payment/common"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Check(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)
	var err error

	req := new(CheckRequest)

	err = utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 查询所有数据
	recordInDb, errorCode := fetchRecord(ctx, req)
	if errorCode != conf.Success {
		modules.BaseError(ctx, errorCode)
		return
	}

	resp := common.PickupResponse(recordInDb)

	modules.BaseSuccess(ctx, resp)
}

func fetchRecord(ctx *gin.Context, req *CheckRequest) (*api_define.TxnReq, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)
	var err error

	resp := new(api_define.TxnReq)

	// 查找交易记录
	resp.TxnRecord, err = record.TxnRecordDao.GetByIDOrUuid(req.MerchantId, req.TxnID, req.PartnerUUID)
	if err != nil {
		logger.Error("GetByIDOrUuid fail->", err.Error())
		return nil, conf.DBError
	}
	if resp.TxnRecord == nil {
		logger.Warn("GetByIDOrUuid no record")
		return nil, conf.RecordNotFund
	}

	// 查找detail
	resp.TxnRecordDetail, err = record.TxnRecordDetailDao.Get(resp.TxnRecord.ID)
	if err != nil {
		logger.Error("recordDetailBean get id->", resp.TxnRecord.ID, ",fail->", err.Error())
		return nil, conf.DBError
	}
	if resp.TxnRecordDetail == nil {
		logger.Warn("recordDetailBean get no record")
		return nil, conf.RecordNotFund
	}

	// 查找merchant account
	if resp.TxnRecord.MerchantAccountID != 0 {
		resp.PaymentProcessRule = new(paymentprocessrule.PaymentProcessRule)
		resp.PaymentProcessRule.MerchantAccount, err = merchantaccount.Dao.
			Get(resp.TxnRecord.MerchantAccountID)
		if err != nil {
			logger.Error("merchantAccount.Get id->", resp.TxnRecord.MerchantAccountID, ",fail->", err.Error())
			return nil, conf.DBError
		}

		if resp.PaymentProcessRule.MerchantAccount == nil {
			logger.Error("merchantAccount.Get id->", resp.TxnRecord.MerchantAccountID, ", no record")
			return nil, conf.UnknownError
		}

		if resp.TxnRecord.TerminalID != 0 {
			resp.PaymentProcessRule.MerchantAccount.Terminal, err = acquirer.TerminalDao.Get(resp.TxnRecord.TerminalID)
			if err != nil {
				logger.Error("terminalBean.Get id->", resp.TxnRecord.TerminalID, ",fail->", err.Error())
				return nil, conf.DBError
			}

			if resp.PaymentProcessRule.MerchantAccount.Terminal == nil {
				logger.Error("merchantAccount.Get id->", resp.TxnRecord.TerminalID, ", no record")
				return nil, conf.UnknownError
			}
		}
	}

	return resp, conf.Success
}
