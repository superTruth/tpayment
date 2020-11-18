package payment

type TxnReq struct {
	Uuid          string `json:"uuid"`
	TxnType       string `json:"txn_type"`
	DeviceID      string `json:"device_id"`
	PaymentMethod string `json:"payment_method"`

	MerchantID string `json:"merchant_id"`

	Amount   string `json:"amount"`
	Currency string `json:"currency"`

	CreditCardBean      *CreditCardBean      `json:"credit_card"`
	CreditCardTokenBean *CreditCardTokenBean `json:"credit_card_token_bean"`
	ApplePayBean        *ApplePayBean        `json:"apple_pay"`
	CreditCard3DSBean   *CreditCard3DSBean   `json:"credit_card_3ds"`
	ConsumerPresentQR   *ConsumerPresentQR   `json:"consumer_present_qr"`

	RealPaymentMethod string `json:"real_payment_method"`
	RealEntryType     string `json:"real_entry_type"`
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
	IsMsdCard               string `json:"is_msd_card"`
	Cvv                     string `json:"cvv"`
	IccRequest              string `json:"icc_request"`
	PIN                     string `json:"pin"`
	ECI                     string `json:"eci"`
	OnlinePaymentCryptogram string `json:"online_payment_cryptogram"`
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
