package payment

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/payment/binddevice"
	"tpayment/models/payment/paymentprocessrule"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func matchProcessRule(ctx *gin.Context, txn *api_define.TxnReq) (*paymentprocessrule.PaymentProcessRule, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)

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

	// 如果有绑定设备匹配，则优先匹配绑定设备Rule
	if txn.DeviceID != "" {
		bindBean := new(binddevice.BindDevice)
		for i, payRule := range matchRules {
			payRule.BindDevice, err = bindBean.Get(models.DB(), ctx,
				payRule.ID, txn.DeviceID)
			if err != nil {
				return nil, conf.DBError
			}
			// 没找到
			if payRule.BindDevice == nil {
				continue
			}
			// 找到
			return payRules[i], conf.SUCCESS
		}
	}

	return matchRules[0], conf.SUCCESS
}
