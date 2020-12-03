package record

import (
	"time"
	"tpayment/models"
)

type TxnRecordDetail struct {
	models.BaseModel

	TxnExpAt *time.Time `gorm:"column:txn_exp_at"`

	CreditCardExp         string `gorm:"column:credit_card_exp"`
	CreditCardFallBack    bool   `gorm:"column:credit_card_fall_back"`
	CreditCardSN          string `gorm:"column:credit_card_sn"`
	CreditCardHolderName  string `gorm:"column:credit_card_holder_name"`
	CreditCardIsMsdCard   bool   `gorm:"column:credit_card_is_msd_card"`
	CreditCardIccRequest  string `gorm:"column:credit_card_icc_request"`
	CreditCardECI         string `gorm:"column:credit_card_eci"`
	CreditCardIccResponse string `gorm:"column:credit_card_icc_response"`
	ResponseCode          string `gorm:"column:response_code"`

	TokenType string `gorm:"column:token_type"`
	Token     string `gorm:"column:token"`

	TDSEnable          bool   `gorm:"column:tds_enable"`
	PayRedirectUrl     string `gorm:"column:pay_redirect_url"`
	RedirectSuccessUrl string `gorm:"column:redirect_success_url"`
	RedirectFailUrl    string `gorm:"column:redirect_fail_url"`

	Addition string `gorm:"column:addition"`
}

func (TxnRecordDetail) TableName() string {
	return "payment_txn_record_detail"
}
