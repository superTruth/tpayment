package void

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/models/payment/acquirer"
	"tpayment/models/payment/merchantaccount"
	"tpayment/models/payment/paymentprocessrule"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func fetchMerchantAccount(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)

	txn.PaymentProcessRule = new(paymentprocessrule.PaymentProcessRule)
	// 查找merchant account
	var err error
	merchantBean := new(merchantaccount.MerchantAccount)
	txn.PaymentProcessRule.MerchantAccount, err =
		merchantBean.Get(models.DB(), ctx, txn.OrgRecord.MerchantAccountID)
	if err != nil {
		logger.Error("merchantBean.Get->", err.Error())
		return conf.DBError
	}
	if txn.PaymentProcessRule.MerchantAccount == nil {
		logger.Warn("can't find merchant account in payment process id->", txn.PaymentProcessRule.ID)
		return conf.ProcessRuleSettingError
	}

	// 查找acquirer
	acquirerBean := new(agency.Acquirer)
	txn.PaymentProcessRule.MerchantAccount.Acquirer, err =
		acquirerBean.Get(models.DB(), ctx, txn.PaymentProcessRule.MerchantAccount.AcquirerID)
	if err != nil {
		logger.Error("acquirerBean.Get->", err.Error())
		return conf.DBError
	}
	if txn.PaymentProcessRule.MerchantAccount.Acquirer == nil {
		logger.Warn("can't find acquirer in merchant id->", txn.PaymentProcessRule.MerchantAccount.ID)
		return conf.ProcessRuleSettingError
	}

	// 查找TID
	if txn.OrgRecord.TerminalID != 0 {
		terminalID := &acquirer.Terminal{
			BaseModel: models.BaseModel{
				Db:  models.DB(),
				Ctx: ctx,
			},
		}
		terminalID, err = terminalID.Get(txn.OrgRecord.TerminalID)
		if err != nil {
			logger.Error("terminalID.Get ", txn.OrgRecord.TerminalID, "->", err.Error())
			return conf.DBError
		}
		if terminalID == nil {
			logger.Error("can't find the terminal ", txn.OrgRecord.TerminalID)
			return conf.RecordNotFund
		}
		terminalID.BaseModel = models.BaseModel{
			Db:  models.DB(),
			Ctx: ctx,
		}
		txn.PaymentProcessRule.MerchantAccount.Terminal = terminalID
	}

	return conf.Success
}
