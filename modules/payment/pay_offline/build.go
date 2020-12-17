package pay_offline

import (
	"tpayment/api/api_define"
	"tpayment/models/payment/record"
)

func preBuildResp(req *api_define.TxnReq) *api_define.TxnResp {
	ret := &api_define.TxnResp{
		Uuid:             req.Uuid,
		TxnType:          req.TxnType,
		DeviceID:         req.DeviceID,
		PaymentMethod:    req.PaymentMethod,
		MerchantID:       req.MerchantID,
		Amount:           req.Amount,
		Currency:         req.Currency,
		TransactionState: record.Success,
	}

	return ret
}

func mergeRespAfterPreHandle(resp *api_define.TxnResp, req *api_define.TxnReq) {
	resp.TxnID = req.TxnRecord.ID
}
