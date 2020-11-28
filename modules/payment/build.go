package payment

import (
	"time"
	"tpayment/api/api_define"
	"tpayment/internal/acquirer_impl"
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
		TransactionState: record.Pending,
	}

	return ret
}

func mergeRespAfterPreHandle(resp *api_define.TxnResp, req *api_define.TxnReq) {
	resp.TxnID = req.TxnRecord.ID
	resp.AcquirerMerchantID = req.PaymentProcessRule.MerchantAccount.MID
	resp.AcquirerName = req.PaymentProcessRule.MerchantAccount.Acquirer.Name
	resp.AcquirerType = req.RealPaymentMethod

	// Terminal
	if req.PaymentProcessRule.MerchantAccount.Terminal != nil &&
		req.PaymentProcessRule.MerchantAccount.Terminal.TID != "" {
		resp.AcquirerTerminalID = req.PaymentProcessRule.MerchantAccount.Terminal.TID
	}
}

func mergeAcquirerResponse(resp *api_define.TxnResp, acquirerResp *acquirer_impl.SaleResponse) {
	resp.AcquirerRRN = acquirerResp.TxnResp.AcquirerRRN
	resp.TransactionState = acquirerResp.TxnResp.TransactionState
	resp.ErrorDesc = acquirerResp.TxnResp.ErrorDesc
	resp.ErrorCode = acquirerResp.TxnResp.ErrorCode
	// 拷贝信用卡数据
	if acquirerResp.TxnResp.CreditCardBean != nil {
		if resp.CreditCardBean == nil {
			resp.CreditCardBean = new(api_define.CreditCardBean)
		}

		resp.CreditCardBean.BatchNum = acquirerResp.TxnResp.CreditCardBean.BatchNum
		resp.CreditCardBean.TraceNum = acquirerResp.TxnResp.CreditCardBean.TraceNum
		resp.CreditCardBean.AuthCode = acquirerResp.TxnResp.CreditCardBean.AuthCode
		resp.CreditCardBean.IccResponse = acquirerResp.TxnResp.CreditCardBean.IccResponse
		resp.CreditCardBean.ResponseCode = acquirerResp.TxnResp.CreditCardBean.ResponseCode
	}
}

func mergeResponseToRecord(record *record.TxnRecord, resp *acquirer_impl.SaleResponse) {
	record.AcquirerRRN = resp.TxnResp.AcquirerRRN
	if resp.TxnResp.CreditCardBean != nil {
		record.AcquirerAuthCode = resp.TxnResp.CreditCardBean.AuthCode
		record.AcquirerBatchNum = resp.TxnResp.CreditCardBean.BatchNum
	}
	record.AcquirerReconID = resp.AcquirerReconID
	t := time.Now()
	record.CompleteAt = &t
	record.Status = resp.TxnResp.TransactionState
}
