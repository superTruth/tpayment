package sale

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/models/payment/acquirer"
	"tpayment/models/payment/merchantaccount"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func fetchMerchantAccount(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
	var errorCode conf.ResultCode

	// 查找merchant account
	var err error
	merchantBean := new(merchantaccount.MerchantAccount)
	txn.PaymentProcessRule.MerchantAccount, err =
		merchantBean.Get(models.DB(), ctx, txn.PaymentProcessRule.MerchantAccountID)
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

	// 分配TID
	if txn.DeviceID != "" {
		tid := &acquirer.Terminal{
			BaseModel: models.BaseModel{
				Db:  models.DB(),
				Ctx: ctx,
			},
		}

		count, err := tid.GetTotal(txn.PaymentProcessRule.MerchantAccountID)
		if err != nil {
			logger.Warn("get total error->", err.Error())
			return conf.DBError
		}
		if count > 0 { // 没有tid，则认为是不需要绑定TID
			tid, errorCode =
				tid.GetOneAvailable(txn.PaymentProcessRule.MerchantAccountID, txn.DeviceID)
			if errorCode != conf.Success {
				logger.Warn("can't get available tid->", txn.DeviceID)
				return errorCode
			}
			txn.PaymentProcessRule.MerchantAccount.Terminal = tid // 可能会没有
		}
	}

	return conf.Success
}
