package payment

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/payment/acquirer"
	"tpayment/models/payment/paymentprocessrule"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func matchProcessRule(ctx *gin.Context, txn *api_define.TxnReq) (*paymentprocessrule.PaymentProcessRule, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)
	var errorCode conf.ResultCode

	rule := new(paymentprocessrule.PaymentProcessRule)

	payRules, err := rule.GetByMerchantID(models.DB(), ctx, txn.MerchantID)
	if err != nil {
		logger.Error("GetByMerchantID fail->", err.Error())
		return nil, conf.DBError
	}

	// 筛选出所有匹配支付方式，支付类型，用卡方式的payment rule
	var matchRules []*paymentprocessrule.PaymentProcessRule
	for i, payRule := range payRules {
		match := false
		// 匹配payment method
		if payRule.PaymentMethods != nil && len(*payRule.PaymentMethods) != 0 {
			for _, method := range *payRule.PaymentMethods {
				if method == txn.RealPaymentMethod {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		// 匹配entry type
		match = false
		if payRule.PaymentEntryTypes != nil && len(*payRule.PaymentEntryTypes) != 0 {
			for _, entryType := range *payRule.PaymentEntryTypes {
				if entryType == txn.RealEntryType {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		// 匹配支付方式
		match = false
		if payRule.PaymentTypes != nil && len(*payRule.PaymentTypes) != 0 {
			for _, paymentTypes := range *payRule.PaymentTypes {
				if paymentTypes == txn.TxnType {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		matchRules = append(matchRules, payRules[i])
	}

	// 未匹配到任何rule
	if len(matchRules) == 0 {
		return nil, conf.NoPaymentProcessRule
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
			return nil, conf.DBError
		}
		if count > 0 { // 没有tid，则认为是不需要绑定TID
			tid, errorCode =
				tid.GetOneAvailable(txn.PaymentProcessRule.MerchantAccountID, txn.DeviceID)
			if errorCode != conf.Success {
				logger.Warn("can't get available tid->", txn.DeviceID)
				return nil, errorCode
			}
			txn.PaymentProcessRule.MerchantAccount.TID = tid // 可能会没有
		}
	}

	return matchRules[0], conf.Success
}
