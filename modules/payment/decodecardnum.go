package payment

import "tpayment/conf"

func analysisPaymentMethod(txn *TxnReq) (string, conf.ResultCode) {
	switch txn.PaymentMethod {
	case RequestCreditCard:
	}
	return "", conf.NotSupport
}
