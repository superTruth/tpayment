package merchantaccount

import (
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/models/payment/acquirer"

	"github.com/gin-gonic/gin"
)

type MerchantAccount struct {
	models.BaseModel
	Name       string              `gorm:"column:name"`
	Currencies *models.StringArray `gorm:"column:currencies;type:JSON"`
	AcquirerID uint                `gorm:"column:acquirer_id"`
	MID        string              `gorm:"column:mid"`
	Addition   string              `gorm:"column:addition"`
	Disable    bool                `gorm:"column:disable"`

	Acquirer *agency.Acquirer   `gorm:"-"`
	Terminal *acquirer.Terminal `gorm:"-"`
}

func (MerchantAccount) TableName() string {
	return "payment_merchant_account"
}

func (m *MerchantAccount) Get(db *models.MyDB, ctx *gin.Context, merchantID uint) (*MerchantAccount, error) {
	var ret = new(MerchantAccount)
	err := db.Model(m).Where("id=?", merchantID).Find(ret).Error
	if err != nil {
		return nil, err
	}

	return ret, nil
}
