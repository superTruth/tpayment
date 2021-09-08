package record

import (
	"database/sql/driver"
	"time"
	"tpayment/models"
	"tpayment/pkg/gorm_json"

	"gorm.io/gorm"
)

var TxnRecordDetailDao = &TxnRecordDetail{}

type TxnRecordDetail struct {
	models.BaseModel

	TxnExpAt *time.Time `gorm:"column:txn_exp_at"`

	CreditCardExp             string `gorm:"column:credit_card_exp"`
	CreditCardFallBack        bool   `gorm:"column:credit_card_fall_back"`
	CreditCardSN              string `gorm:"column:credit_card_sn"`
	CreditCardHolderName      string `gorm:"column:credit_card_holder_name"`
	CreditCardIsMsdCard       bool   `gorm:"column:credit_card_is_msd_card"`
	CreditCardIccRequest      string `gorm:"column:credit_card_icc_request"`
	CreditCardIccResponse     string `gorm:"column:credit_card_icc_response"`
	CreditCardIccScriptUpload string `gorm:"column:credit_card_icc_script_upload"`
	CreditCardECI             string `gorm:"column:credit_card_eci"`

	ResponseCode string `gorm:"column:response_code"`

	TokenType string `gorm:"column:token_type"`
	Token     string `gorm:"column:token"`

	TDSEnable          bool   `gorm:"column:tds_enable"`
	PayRedirectUrl     string `gorm:"column:pay_redirect_url"`
	RedirectSuccessUrl string `gorm:"column:redirect_success_url"`
	RedirectFailUrl    string `gorm:"column:redirect_fail_url"`
	ResultNotifyUrl    string `gorm:"column:result_notify_url"`
	Signature          string `gorm:"column:signature"`

	Addition *AdditionData `gorm:"column:addition;type:JSON"`
}

func (TxnRecordDetail) TableName() string {
	return "payment_txn_record_detail"
}

func (t *TxnRecordDetail) Create(detail *TxnRecordDetail) error {
	return models.DB.Model(t).Create(detail).Error
}

func (t *TxnRecordDetail) Get(id uint64) (*TxnRecordDetail, error) {
	ret := new(TxnRecordDetail)
	err := models.DB.Model(t).Where("id=?", id).First(ret).Error
	if gorm.ErrRecordNotFound == err { // 没有记录, 就创建一条记录
		if err != nil {
			return nil, err
		}
		return t, err
	}
	return ret, err
}

func (t *TxnRecordDetail) Update(db *gorm.DB) error {
	return db.Model(t).Save(t).Error
}

type AdditionData struct {
	CupTraceNum string `json:"cup_trace_num" gorm:"column:cup_trace_num"`
	CupRrn      string `json:"cup_rrn" gorm:"column:cup_rrn"`
}

//
func (c AdditionData) Value() (driver.Value, error) {
	return gorm_json.Value(c)
}
func (c *AdditionData) Scan(input interface{}) error {
	return gorm_json.Scan(input, c)
}
