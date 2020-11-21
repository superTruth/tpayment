package merchantaccount

import (
	"tpayment/models"
	"tpayment/models/agency"

	"github.com/gin-gonic/gin"
)

type MerchantAccount struct {
	models.BaseModel
	Name       string              `json:"name"`
	Currencies *models.StringArray `json:"currencies"`
	AcquirerID uint                `json:"acquirer_id"`
	MID        string              `json:"mid"`
	Addition   string              `json:"addition"`
	Disable    bool                `json:"disable"`

	Acquirer *agency.Acquirer `json:"-"`
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
