package paymentprocessrule

import (
	"tpayment/models"
	"tpayment/models/payment/binddevice"
	"tpayment/models/payment/merchantaccount"

	"github.com/gin-gonic/gin"
)

type PaymentProcessRule struct {
	models.BaseModel

	MerchantID uint `json:"merchant_id"`

	MerchantAccountID uint                `json:"merchant_account_id"`
	PaymentMethods    *models.StringArray `json:"payment_methods"`
	PaymentEntryTypes *models.StringArray `json:"payment_entry_types"`
	PaymentTypes      *models.StringArray `json:"payment_types"`

	BindDevice      *binddevice.BindDevice           `json:"-"`
	MerchantAccount *merchantaccount.MerchantAccount `json:"-"`
}

func (PaymentProcessRule) TableName() string {
	return "payment_process_rule"
}

func (rule *PaymentProcessRule) GetByMerchantID(db *models.MyDB, ctx *gin.Context, merchantID uint) ([]*PaymentProcessRule, error) {
	var ret []*PaymentProcessRule
	err := db.Model(rule).Where("merchant_id=?", merchantID).Find(ret).Error
	if err != nil {
		return nil, err
	}

	return ret, nil
}
