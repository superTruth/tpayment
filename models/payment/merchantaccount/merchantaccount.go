package merchantaccount

import (
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/models/payment/acquirer"

	"github.com/gin-gonic/gin"
)

type MerchantAccount struct {
	models.BaseModel
	Name       string              `gorm:"name"`
	Currencies *models.StringArray `gorm:"currencies"`
	AcquirerID uint                `gorm:"acquirer_id"`
	MID        string              `gorm:"mid"`
	Addition   string              `gorm:"addition"`
	Disable    bool                `gorm:"disable"`

	Acquirer *agency.Acquirer   `gorm:"-"`
	Terminal *acquirer.Terminal `gorm:"-"`
}

func (MerchantAccount) TableName() string {
	return "merchant_account"
}

func (m *MerchantAccount) Get(db *models.MyDB, ctx *gin.Context, merchantID uint) (*MerchantAccount, error) {
	var ret *MerchantAccount
	err := db.Model(m).Where("id=?", merchantID).Find(ret).Error
	if err != nil {
		return nil, err
	}

	return ret, nil
}
