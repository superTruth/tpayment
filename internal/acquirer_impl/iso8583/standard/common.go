package standard

import "tpayment/internal/acquirer_impl"

func GetAccountTag(req *acquirer_impl.SaleRequest) string {
	var mid, tid string
	mid = req.TxqReq.PaymentProcessRule.MerchantAccount.MID

	if req.TxqReq.PaymentProcessRule.BindDevice != nil {
		tid = req.TxqReq.PaymentProcessRule.BindDevice.TID
	}

	return mid + tid
}

func RandomTID(req *acquirer_impl.SaleRequest) error {
	// 已经指定了TID，就不需要随机TID
	if req.TxqReq.PaymentProcessRule.BindDevice != nil {
		return nil
	}

	return nil
}

// TODO need
func LockTID() error {
	return nil
}

func UnlockTID() error {
	return nil
}
