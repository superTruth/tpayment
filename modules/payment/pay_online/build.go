package pay_online

import (
	"fmt"
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
	if req.PaymentProcessRule.MerchantAccount.Terminal != nil {
		resp.AcquirerTerminalID = req.PaymentProcessRule.MerchantAccount.Terminal.TID

		if resp.CreditCardBean == nil {
			resp.CreditCardBean = new(api_define.CreditCardBean)
		}
		resp.CreditCardBean.BatchNum = req.CreditCardBean.BatchNum
		resp.CreditCardBean.TraceNum = req.CreditCardBean.TraceNum
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

		resp.CreditCardBean.AuthCode = acquirerResp.TxnResp.CreditCardBean.AuthCode
		resp.CreditCardBean.IccResponse = acquirerResp.TxnResp.CreditCardBean.IccResponse
		resp.CreditCardBean.ResponseCode = acquirerResp.TxnResp.CreditCardBean.ResponseCode
	}
}

func mergeResponseToRecord(req *api_define.TxnReq, resp *acquirer_impl.SaleResponse) {
	req.TxnRecord.AcquirerRRN = resp.TxnResp.AcquirerRRN

	req.TxnRecord.AcquirerReconID = resp.AcquirerReconID
	t := time.Now()
	req.TxnRecord.CompleteAt = &t
	req.TxnRecord.AcquirerTxnDateTime = resp.TxnResp.DateTime
	fmt.Println("UpdateTxnResult2->", req.TxnRecord.AcquirerTxnDateTime)
	req.TxnRecord.Status = resp.TxnResp.TransactionState
	req.TxnRecord.ErrorCode = resp.TxnResp.ErrorCode
	req.TxnRecord.ErrorDes = resp.TxnResp.ErrorDesc

	if resp.TxnResp.CreditCardBean != nil {
		req.TxnRecordDetail.CreditCardIccResponse = resp.TxnResp.CreditCardBean.IccResponse
		req.TxnRecordDetail.ResponseCode = resp.TxnResp.CreditCardBean.ResponseCode
	}
	req.TxnRecordDetail.PayRedirectUrl = resp.TxnResp.PayRedirectUrl
	req.TxnRecordDetail.Addition = resp.TxnResp.AdditionData
}
