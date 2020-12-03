package paymentprocessrule

import (
	"tpayment/models"
	"tpayment/models/payment/merchantaccount"

	"github.com/gin-gonic/gin"
)

type PaymentProcessRule struct {
	models.BaseModel

	MerchantID uint64 `gorm:"column:merchant_id"`

	MerchantAccountID uint64              `gorm:"column:merchant_account_id"`
	PaymentMethods    *models.StringArray `gorm:"column:payment_methods;type:JSON"`
	PaymentEntryTypes *models.StringArray `gorm:"column:payment_entry_types;type:JSON"`
	PaymentTypes      *models.StringArray `gorm:"column:payment_types;type:JSON"`

	MerchantAccount *merchantaccount.MerchantAccount `json:"-"`
}

func (PaymentProcessRule) TableName() string {
	return "payment_process_rule"
}

func (rule *PaymentProcessRule) GetByMerchantID(db *models.MyDB, ctx *gin.Context, merchantID uint64) ([]*PaymentProcessRule, error) {
	var ret []*PaymentProcessRule
	err := db.Model(rule).Where("merchant_id=?", merchantID).Find(&ret).Error
	if err != nil {
		return nil, err
	}
	return ret, nil
}
