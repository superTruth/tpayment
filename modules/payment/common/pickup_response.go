package common

import (
	"tpayment/api/api_define"
	"tpayment/models/payment/record"
)

func PickupResponse(req *api_define.TxnReq) *api_define.TxnResp {
	ret := &api_define.TxnResp{
		Uuid:             req.TxnRecord.PartnerUUID,
		TxnID:            req.TxnRecord.ID,
		TxnType:          req.TxnRecord.PaymentType,
		DeviceID:         req.TxnRecord.PaymentFromDeviceID,
		PaymentMethod:    req.TxnRecord.PaymentMethod,
		MerchantID:       req.TxnRecord.MerchantID,
		TotalAmount:      req.TxnRecord.TotalAmount.String(),
		Amount:           req.TxnRecord.Amount.String(),
		Currency:         req.TxnRecord.Currency,
		TransactionState: req.TxnRecord.Status,
		ErrorCode:        req.TxnRecord.ErrorCode,
		ErrorDesc:        req.TxnRecord.ErrorDes,
		DateTime:         req.TxnRecord.CompleteAt,
		AcquirerRRN:      req.TxnRecord.AcquirerRRN,
		AcquirerType:     req.TxnRecord.PaymentMethod,
		PayRedirectUrl:   req.TxnRecordDetail.PayRedirectUrl,
		AdditionData:     req.TxnRecordDetail.Addition,
	}

	if req.PaymentProcessRule != nil && req.PaymentProcessRule.MerchantAccount != nil {
		ret.AcquirerMerchantID = req.PaymentProcessRule.MerchantAccount.MID
		if req.PaymentProcessRule.MerchantAccount.Acquirer != nil {
			ret.AcquirerName = req.PaymentProcessRule.MerchantAccount.Acquirer.Name
		}
		if req.PaymentProcessRule.MerchantAccount.Terminal != nil {
			ret.AcquirerTerminalID = req.PaymentProcessRule.MerchantAccount.Terminal.TID
		}
	}

	if req.TxnRecordDetail.CreditCardExp != "" {
		ret.CreditCardBean = &api_define.CreditCardBean{
			CardReaderMode: req.TxnRecord.PaymentEntryType,
			CardExpMonth:   req.TxnRecordDetail.CreditCardExp[2:],
			CardExpYear:    req.TxnRecordDetail.CreditCardExp[:2],
			CardFallback:   req.TxnRecordDetail.CreditCardFallBack,
			CardNumber:     req.TxnRecord.ConsumerIdentify,
			CardSn:         req.TxnRecordDetail.CreditCardSN,
			CardHolderName: req.TxnRecordDetail.CreditCardHolderName,
			IsMsdCard:      req.TxnRecordDetail.CreditCardIsMsdCard,
			IccRequest:     req.TxnRecordDetail.CreditCardIccRequest,
			IccResponse:    req.TxnRecordDetail.CreditCardIccResponse,
			TraceNum:       req.TxnRecord.AcquirerTraceNum,
			BatchNum:       req.TxnRecord.AcquirerBatchNum,
			AuthCode:       req.TxnRecord.AcquirerAuthCode,
			ResponseCode:   req.TxnRecordDetail.ResponseCode,
		}
	}

	// 以下几种状态需要转换一下
	switch ret.TransactionState {
	case record.Init:
		ret.TransactionState = record.Pending
	case record.NeedReversal, record.Reversal:
		ret.TransactionState = record.Fail
	}

	return ret
}
