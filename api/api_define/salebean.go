package api_define

import (
	"errors"
	"time"
	"tpayment/conf"
	"tpayment/models/agency"
	"tpayment/models/merchant"
	"tpayment/models/payment/paymentprocessrule"
	"tpayment/models/payment/record"
	"tpayment/pkg/paymentmethod/decodecardnum/creditcard"

	"github.com/gin-gonic/gin"
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
	DateTime              *time.Time `json:"date_time,omitempty"`
	AcquirerMerchantID    string     `json:"acquirer_merchant_id"`
	AcquirerTerminalID    string     `json:"acquirer_terminal_id"`
	AcquirerRRN           string     `json:"acquirer_rrn,omitempty"`
	AcquirerName          string     `json:"acquirer_name,omitempty"`
	AcquirerType          string     `json:"acquirer_type,omitempty"`
	AcquirerReconID       string     `json:"acquirer_recon_id"`
	InvoiceNum            uint64     `json:"invoice_num"`
	CustomerPaymentMethod string     `json:"customer_payment_method"`

	CreditCardBean      *CreditCardBean      `json:"credit_card,omitempty"`
	CreditCardTokenBean *CreditCardTokenBean `json:"credit_card_token_bean,omitempty"`
	ApplePayBean        *ApplePayBean        `json:"apple_pay,omitempty"`
	CreditCard3DSBean   *CreditCard3DSBean   `json:"credit_card_3ds,omitempty"`
	ConsumerPresentQR   *ConsumerPresentQR   `json:"consumer_present_qr,omitempty"`
	AdditionData        *record.AdditionData `json:"addition_data,omitempty"`

	RedirectSuccessUrl string `json:"redirect_success_url,omitempty"`
	RedirectFailUrl    string `json:"redirect_fail_url,omitempty"`
	ResultNotifyUrl    string `json:"result_notify_url,omitempty"`

	TxnExpAt *time.Time `json:"txn_exp_at"`

	// 用于void,refund，tips，等二次交易
	OrgTxnID uint64 `json:"origin_txn_id,omitempty"`

	// 后期处理填充
	RealPaymentMethod  string                                 `json:"real_payment_method,omitempty"`
	RealEntryType      string                                 `json:"real_entry_type,omitempty"`
	PaymentProcessRule *paymentprocessrule.PaymentProcessRule `json:"payment_process_rule,omitempty"`
	FromName           string                                 `json:"from_name,omitempty"`
	FromIp             string                                 `json:"from_ip,omitempty"`
	CashierID          string                                 `json:"cashier_id,omitempty"`
	TxnRecord          *record.TxnRecord                      `json:"txn_record,omitempty"`
	TxnRecordDetail    *record.TxnRecordDetail                `json:"txn_record_detail,omitempty"`
	OrgRecord          *record.TxnRecord                      `json:"org_record,omitempty"`
	OrgRecordDetail    *record.TxnRecordDetail                `json:"org_record_detail,omitempty"`
	MerchantInfo       *merchant.Merchant                     `json:"merchant_info,omitempty"`
	AgencyInfo         *agency.Agency                         `json:"agency_info,omitempty"`
}

type CreditCardBean struct {
	CardReaderMode          string `json:"card_reader_mode,omitempty"`
	CardExpMonth            string `json:"card_exp_month,omitempty"`
	CardExpYear             string `json:"card_exp_year,omitempty"`
	CardFallback            bool   `json:"card_fallback,omitempty"`
	CardNumber              string `json:"card_number,omitempty"`
	CardSn                  string `json:"card_sn,omitempty"`
	CardTrack1              string `json:"card_track1,omitempty"`
	CardTrack2              string `json:"card_track2,omitempty"`
	CardTrack3              string `json:"card_track3,omitempty"`
	CardHolderName          string `json:"card_holder_name,omitempty"`
	IsMsdCard               bool   `json:"is_msd_card,omitempty"`
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

type Transfer struct {
	DestAccount string `json:"dest_account"`
}

// 回复数据
type TxnResp struct {
	Uuid          string `json:"uuid,omitempty"`
	TxnID         uint64 `json:"txn_id,omitempty"`
	TxnType       string `json:"txn_type,omitempty"`
	DeviceID      string `json:"device_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`

	MerchantID uint64 `json:"merchant_id,omitempty"`

	TotalAmount string `json:"total_amount,omitempty"`
	Amount      string `json:"amount,omitempty"`
	Currency    string `json:"currency,omitempty"`

	TransactionState string `json:"transaction_state,omitempty"`
	ErrorCode        string `json:"error_code,omitempty"`
	ErrorDesc        string `json:"error_desc,omitempty"`

	DateTime           *time.Time `json:"date_time,omitempty"`
	AcquirerMerchantID string     `json:"acquirer_merchant_id,omitempty"`
	AcquirerTerminalID string     `json:"acquirer_terminal_id,omitempty"`
	AcquirerRRN        string     `json:"acquirer_rrn,omitempty"`
	AcquirerName       string     `json:"acquirer_name,omitempty"`
	AcquirerType       string     `json:"acquirer_type,omitempty"`
	InvoiceNum         uint64     `json:"invoice_num"`

	CreditCardBean *CreditCardBean `json:"credit_card,omitempty"`

	PayRedirectUrl string `json:"pay_redirect_url,omitempty"`

	AdditionData *record.AdditionData `json:"addition_data,omitempty"`
}

// 验证请求参数是否正确
func Validate(ctx *gin.Context, txn *TxnReq) error {
	if txn.Uuid == "" {
		return errors.New("uuid is empty")
	}

	if _, ok := conf.PaymentType[txn.TxnType]; !ok {
		return errors.New("can't support txn type->" + txn.TxnType)
	}

	if txn.PaymentMethod != "" {
		if _, ok := conf.RequestPaymentMethod[txn.PaymentMethod]; !ok {
			return errors.New("can't support payment method->" + txn.PaymentMethod)
		}
	}

	if txn.MerchantID == 0 {
		return errors.New("merchant id is empty")
	}

	if txn.Amount != "" {
		if txn.Currency == "" {
			return errors.New("currency is empty")
		}
		if _, ok := conf.CurrencyCode[txn.Currency]; !ok {
			return errors.New("can't support currency->" + txn.PaymentMethod)
		}
	}

	switch txn.PaymentMethod {
	case conf.RequestCreditCard:
		if txn.CreditCardBean == nil {
			return errors.New("credit card information empty")
		}

		if _, ok := conf.PaymentEntryType[txn.CreditCardBean.CardReaderMode]; !ok {
			return errors.New("can't support card read mode->" + txn.CreditCardBean.CardReaderMode)
		}

		// 检查卡片有效期
		ok, err := creditcard.CheckCardExp(txn.CreditCardBean.CardExpYear,
			txn.CreditCardBean.CardExpMonth)
		if err != nil {
			return errors.New("credit card exp format error")
		}
		if !ok {
			return errors.New("card expired")
		}
		if txn.CreditCardBean.CardNumber == "" {
			return errors.New("card number is empty")
		}
	case conf.Transfer:

	case conf.RequestCreditCardToken:
		if txn.CreditCardTokenBean == nil {
			return errors.New("credit card token empty")
		}
		if txn.CreditCardTokenBean.Token == "" {
			return errors.New("credit card token empty")
		}
	case conf.RequestConsumerPresentQR:
		if txn.ConsumerPresentQR == nil {
			return errors.New("consumer present qr empty")
		}
		if txn.ConsumerPresentQR.Content == "" {
			return errors.New("consumer present qr empty")
		}
	case conf.RequestApplePay:
		if txn.ApplePayBean == nil {
			return errors.New("apple pay empty")
		}
		if txn.ApplePayBean.Token == "" {
			return errors.New("apple pay empty")
		}
	case conf.RequestOther:
		if txn.CustomerPaymentMethod == "" {
			return errors.New("customer payment method is empty")
		}
	}

	return nil
}
