package standard

import "tpayment/internal/acquirer_impl"

func GetAccountTag(req *acquirer_impl.SaleRequest) string {
	var mid, tid string
	mid = req.TxqReq.PaymentProcessRule.MerchantAccount.MID

	if req.TxqReq.PaymentProcessRule.MerchantAccount.TID != nil {
		tid = req.TxqReq.PaymentProcessRule.MerchantAccount.TID.TID
	}

	return mid + tid
}
