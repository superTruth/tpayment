package sale

import (
	"errors"
	"time"
	"tpayment/api/api_define"
	"tpayment/internal/acquirer_impl"
	"tpayment/models/payment/record"

	"github.com/shopspring/decimal"
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
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, errors.New("can't parse amount->" + req.Amount + "," + err.Error())
	}

	ret := &record.TxnRecord{
		MerchantID:          req.MerchantID,
		TotalAmount:         amount,
		Amount:              amount,
		Currency:            req.Currency,
		MerchantAccountID:   req.PaymentProcessRule.MerchantID,
		PaymentMethod:       req.RealPaymentMethod,
		PaymentEntryType:    req.RealEntryType,
		PaymentType:         req.TxnType,
		PartnerUUID:         req.Uuid,
		Status:              record.Init,
		PaymentFromName:     req.FromName,
		PaymentFromIP:       req.FromIp,
		PaymentFromDeviceID: req.DeviceID,
		InvoiceNum:          "",
		CashierID:           req.CashierID,
	}

	if req.CreditCardBean != nil && req.CreditCardBean.CardNumber != "" {
		ret.ConsumerIdentify = req.CreditCardBean.CardNumber
	}

	return ret, nil
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
