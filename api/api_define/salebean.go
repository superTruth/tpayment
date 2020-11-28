package api_define

import (
	"tpayment/models/agency"
	"tpayment/models/merchant"
	"tpayment/models/payment/paymentprocessrule"
	"tpayment/models/payment/record"
)

// 请求数据
type TxnReq struct {
	Uuid          string `json:"uuid"`
	TxnType       string `json:"txn_type"`
	DeviceID      string `json:"device_id"`
	PaymentMethod string `json:"payment_method"`

	MerchantID uint `json:"merchant_id"`

	Amount   string `json:"amount"`
	Currency string `json:"currency"`

	CreditCardBean      *CreditCardBean      `json:"credit_card"`
	CreditCardTokenBean *CreditCardTokenBean `json:"credit_card_token_bean"`
	ApplePayBean        *ApplePayBean        `json:"apple_pay"`
	CreditCard3DSBean   *CreditCard3DSBean   `json:"credit_card_3ds"`
	ConsumerPresentQR   *ConsumerPresentQR   `json:"consumer_present_qr"`

	// 用于void,refund，tips，等二次交易
	OriginTxnID uint `json:"origin_txn_id"`

	// 后期处理填充
	RealPaymentMethod  string                                 `json:"real_payment_method"`
	RealEntryType      string                                 `json:"real_entry_type"`
	PaymentProcessRule *paymentprocessrule.PaymentProcessRule `json:"payment_process_rule"`
	FromName           string                                 `json:"from_name"`
	FromIp             string                                 `json:"from_ip"`
	CashierID          string                                 `json:"cashier_id"`
	TxnRecord          *record.TxnRecord                      `json:"txn_record"`
	OrgRecord          *record.TxnRecord                      `json:"org_record"`
	MerchantInfo       *merchant.Merchant                     `json:"merchant_info"`
	AgencyInfo         *agency.Agency                         `json:"agency_info"`
}

type CreditCardBean struct {
	CardReaderMode          string `json:"card_reader_mode"`
	CardExpMonth            string `json:"card_exp_month"`
	CardExpYear             string `json:"card_exp_year"`
	CardExpDay              string `json:"card_exp_day"`
	CardFallback            bool   `json:"card_fallback"`
	CardNumber              string `json:"card_number"`
	CardSn                  string `json:"card_sn"`
	CardTrack1              string `json:"card_track1"`
	CardTrack2              string `json:"card_track2"`
	CardTrack3              string `json:"card_track3"`
	CardHolderName          string `json:"card_holder_name"`
	IsMsdCard               bool   `json:"is_msd_card"`
	Cvv                     string `json:"cvv"`
	IccRequest              string `json:"icc_request"`
	PIN                     string `json:"pin"`
	ECI                     string `json:"eci"`
	OnlinePaymentCryptogram string `json:"online_payment_cryptogram"`

	IccResponse  string `json:"icc_response"`
	TraceNum     uint   `json:"trace_num"`
	BatchNum     uint   `json:"batch_num"`
	AuthCode     string `json:"auth_code"`
	ResponseCode string `json:"response_code"`
}

type CreditCardTokenBean struct {
	Token string `json:"token"`
}

type ApplePayBean struct {
	Token string `json:"token"`
}

type CreditCard3DSBean struct {
	Enable     bool   `json:"enable"`
	SuccessUrl string `json:"success_url"`
	FailUrl    string `json:"fail_url"`
}

type ConsumerPresentQR struct {
	CodeType string `json:"code_type"`
	Content  string `json:"content"`
}

// 回复数据
type TxnResp struct {
	Uuid          string `json:"uuid"`
	TxnID         uint   `json:"txn_id"`
	TxnType       string `json:"txn_type"`
	DeviceID      string `json:"device_id"`
	PaymentMethod string `json:"payment_method"`

	MerchantID uint `json:"merchant_id"`

	Amount   string `json:"amount"`
	Currency string `json:"currency"`

	TransactionState string `json:"transaction_state"`
	ErrorCode        string `json:"error_code"`
	ErrorDesc        string `json:"error_desc"`

	DateTime           string `json:"date_time"`
	AcquirerMerchantID string `json:"acquirer_merchant_id"`
	AcquirerTerminalID string `json:"acquirer_terminal_id"`
	AcquirerRRN        string `json:"acquirer_rrn"`
	AcquirerName       string `json:"acquirer_name"`
	AcquirerType       string `json:"acquirer_type"`

	CreditCardBean *CreditCardBean `json:"credit_card"`
}
