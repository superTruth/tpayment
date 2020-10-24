package payment

type TxnReq struct {
	Uuid          string `json:"uuid"`
	TxnType       string `json:"txn_type"`
	DeviceID      string `json:"device_id"`
	PaymentMethod string `json:"payment_method"`

	Amount   string `json:"amount"`
	Currency string `json:"currency"`

	CreditCardBean *CreditCardBean `json:"credit_card"`
}

type CreditCardBean struct {
	CardReaderMode string `json:"card_reader_mode"`
	CardExpMonth   string `json:"card_exp_month"`
	CardExpYear    string `json:"card_exp_year"`
	CardFallback   bool   `json:"card_fallback"`
	CardNumber     string `json:"card_number"`
	CardSn         string `json:"card_sn"`
	CardTrack1     string `json:"card_track1"`
	CardTrack2     string `json:"card_track2"`
	CardTrack3     string `json:"card_track3"`
	CardHolderName string `json:"card_holder_name"`
	IsMsdCard      string `json:"is_msd_card"`
	Cvv            string `json:"cvv"`
	IccRequest     string `json:"icc_request"`
	PIN            string `json:"pin"`
	Token          string `json:"token"`
}
