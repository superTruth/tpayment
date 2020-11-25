package payment

import (
	"tpayment/api/api_define"
	"tpayment/models/payment/record"
)

func preBuildResp(req *api_define.TxnReq) *api_define.TxnResp {
	ret := &api_define.TxnResp{
		Uuid:               req.Uuid,
		TxnID:              req.TxnRecord.ID,
		TxnType:            req.TxnType,
		DeviceID:           req.DeviceID,
		PaymentMethod:      req.PaymentMethod,
		MerchantID:         req.MerchantID,
		Amount:             req.Amount,
		Currency:           req.Currency,
		TransactionState:   record.Pending,
		AcquirerMerchantID: req.PaymentProcessRule.MerchantAccount.MID,
		AcquirerTerminalID: "",
		AcquirerRRN:        "",
		AcquirerName:       req.PaymentProcessRule.MerchantAccount.Acquirer.Name,
		AcquirerType:       req.RealPaymentMethod,
		CreditCardBean:     nil,
	}

	// TID

	return ret
}
