package void

import (
	"time"
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
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
		AcquirerName:       req.PaymentProcessRule.MerchantAccount.Acquirer.Name,
		AcquirerType:       req.RealPaymentMethod,
		CreditCardBean:     nil,
	}

	// Terminal
	if req.PaymentProcessRule.MerchantAccount.Terminal != nil &&
		req.PaymentProcessRule.MerchantAccount.Terminal.TID != "" {
		ret.AcquirerTerminalID = req.PaymentProcessRule.MerchantAccount.Terminal.TID
	}

	return ret
}

func mergeAcquirerResponse(resp *api_define.TxnResp, acquirerResp *acquirer_impl.SaleResponse) {
	resp.AcquirerRRN = acquirerResp.TxnResp.AcquirerRRN

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

// 创建历史记录
func buildRecord(req *api_define.TxnReq) (*record.TxnRecord, error) {
	ret := &record.TxnRecord{
		MerchantID:          req.OrgRecord.MerchantID,
		TotalAmount:         req.OrgRecord.Amount,
		Amount:              req.OrgRecord.Amount,
		Currency:            req.OrgRecord.Currency,
		MerchantAccountID:   req.OrgRecord.MerchantAccountID,
		PaymentMethod:       req.OrgRecord.PaymentMethod,
		PaymentEntryType:    conf.ManualInput,
		PaymentType:         req.TxnType,
		PartnerUUID:         req.Uuid,
		Status:              record.Init,
		PaymentFromName:     req.FromName,
		PaymentFromIP:       req.FromIp,
		PaymentFromDeviceID: req.DeviceID,
		InvoiceNum:          "",
		CashierID:           req.CashierID,
		ConsumerIdentify:    req.OrgRecord.ConsumerIdentify,
	}

	return ret, nil
}

func mergeResponseToRecord(req *api_define.TxnReq, resp *acquirer_impl.SaleResponse) {
	req.TxnRecord.AcquirerRRN = resp.TxnResp.AcquirerRRN
	if resp.TxnResp.CreditCardBean != nil {
		req.TxnRecord.AcquirerAuthCode = resp.TxnResp.CreditCardBean.AuthCode
		req.TxnRecord.AcquirerBatchNum = resp.TxnResp.CreditCardBean.BatchNum
	}
	req.TxnRecord.AcquirerReconID = resp.AcquirerReconID
	t := time.Now()
	req.TxnRecord.CompleteAt = &t
	req.TxnRecord.Status = resp.TxnResp.TransactionState
}
