package api_define

import (
	"time"
	"tpayment/models/agency"
	"tpayment/models/merchant"
	"tpayment/models/payment/paymentprocessrule"
	"tpayment/models/payment/record"
)

// 请求数据
type TxnReq struct {
	Uuid          string `json:"uuid,omitempty"`
	TxnType       string `json:"txn_type,omitempty"`
	DeviceID      string `json:"device_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`

	MerchantID uint64 `json:"merchant_id,omitempty"`

	Amount   string `json:"amount"`
	Currency string `json:"currency,omitempty"`

	// 用于offline
	DateTime           *time.Time `json:"date_time,omitempty"`
	AcquirerMerchantID string     `json:"acquirer_merchant_id"`
	AcquirerTerminalID string     `json:"acquirer_terminal_id"`
	AcquirerRRN        string     `json:"acquirer_rrn,omitempty"`
	AcquirerName       string     `json:"acquirer_name,omitempty"`
	AcquirerType       string     `json:"acquirer_type,omitempty"`
	AcquirerReconID    string     `json:"acquirer_recon_id"`

	CreditCardBean      *CreditCardBean      `json:"credit_card,omitempty"`
	CreditCardTokenBean *CreditCardTokenBean `json:"credit_card_token_bean,omitempty"`
	ApplePayBean        *ApplePayBean        `json:"apple_pay,omitempty"`
	CreditCard3DSBean   *CreditCard3DSBean   `json:"credit_card_3ds,omitempty"`
	ConsumerPresentQR   *ConsumerPresentQR   `json:"consumer_present_qr,omitempty"`
	AdditionData        string               `json:"addition_data,omitempty"`

	// 用于void,refund，tips，等二次交易
	OriginTxnID uint64 `json:"origin_txn_id,omitempty"`

	// 后期处理填充
	RealPaymentMethod  string                                 `json:"real_payment_method,omitempty"`
	RealEntryType      string                                 `json:"real_entry_type,omitempty"`
	PaymentProcessRule *paymentprocessrule.PaymentProcessRule `json:"payment_process_rule,omitempty"`
	FromName           string                                 `json:"from_name,omitempty"`
	FromIp             string                                 `json:"from_ip,omitempty"`
	CashierID          string                                 `json:"cashier_id,omitempty"`
	TxnRecord          *record.TxnRecord                      `json:"txn_record,omitempty"`
	TxnRecordDetail    *record.TxnRecordDetail                `json:"txn_record_detail"`
	OrgRecord          *record.TxnRecord                      `json:"org_record,omitempty"`
	MerchantInfo       *merchant.Merchant                     `json:"merchant_info,omitempty"`
	AgencyInfo         *agency.Agency                         `json:"agency_info,omitempty"`
}

type CreditCardBean struct {
	CardReaderMode          string `json:"card_reader_mode,omitempty"`
	CardExpMonth            string `json:"card_exp_month,omitempty"`
	CardExpYear             string `json:"card_exp_year,omitempty"`
	CardExpDay              string `json:"card_exp_day,omitempty"`
	CardFallback            bool   `json:"card_fallback,omitempty"`
	CardNumber              string `json:"card_number,omitempty"`
	CardSn                  string `json:"card_sn,omitempty"`
	CardTrack1              string `json:"card_track1,omitempty"`
	CardTrack2              string `json:"card_track2,omitempty"`
	CardTrack3              string `json:"card_track3,omitempty"`
	CardHolderName          string `json:"card_holder_name,omitempty"`
	IsMsdCard               bool   `json:"is_msd_card,omitempty"`
	IsOffline               bool   `json:"is_offline,omitempty"`
	Cvv                     string `json:"cvv,omitempty"`
	IccRequest              string `json:"icc_request,omitempty"`
	PIN                     string `json:"pin,omitempty"`
	ECI                     string `json:"eci,omitempty"`
	OnlinePaymentCryptogram string `json:"online_payment_cryptogram,omitempty"`

	IccResponse  string `json:"icc_response,omitempty"`
	TraceNum     uint64 `json:"trace_num,omitempty"`
	BatchNum     uint64 `json:"batch_num,omitempty"`
	AuthCode     string `json:"auth_code,omitempty"`
	ResponseCode string `json:"response_code,omitempty"`
}

type CreditCardTokenBean struct {
	Token string `json:"token,omitempty"`
}

type ApplePayBean struct {
	Token string `json:"token,omitempty"`
}

type CreditCard3DSBean struct {
	Enable     bool   `json:"enable,omitempty"`
	SuccessUrl string `json:"success_url,omitempty"`
	FailUrl    string `json:"fail_url,omitempty"`
}

type ConsumerPresentQR struct {
	CodeType string `json:"code_type,omitempty"`
	Content  string `json:"content,omitempty"`
}

// 回复数据
type TxnResp struct {
	Uuid          string `json:"uuid,omitempty"`
	TxnID         uint64 `json:"txn_id,omitempty"`
	TxnType       string `json:"txn_type,omitempty"`
	DeviceID      string `json:"device_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`

	MerchantID uint64 `json:"merchant_id,omitempty"`

	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`

	TransactionState string `json:"transaction_state,omitempty"`
	ErrorCode        string `json:"error_code,omitempty"`
	ErrorDesc        string `json:"error_desc,omitempty"`

	DateTime           string `json:"date_time,omitempty"`
	AcquirerMerchantID string `json:"acquirer_merchant_id,omitempty"`
	AcquirerTerminalID string `json:"acquirer_terminal_id,omitempty"`
	AcquirerRRN        string `json:"acquirer_rrn,omitempty"`
	AcquirerName       string `json:"acquirer_name,omitempty"`
	AcquirerType       string `json:"acquirer_type,omitempty"`

	CreditCardBean *CreditCardBean `json:"credit_card,omitempty"`

	AdditionData string `json:"addition_data,omitempty"`
}
