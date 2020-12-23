package merchantaccount

import (
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/models/payment/acquirer"

	"github.com/jinzhu/gorm"
)

var Dao = &MerchantAccount{}

type MerchantAccount struct {
	models.BaseModel
	Name       string              `gorm:"column:name"`
	Currencies *models.StringArray `gorm:"column:currencies;type:JSON"`
	AcquirerID uint64              `gorm:"column:acquirer_id"`
	MID        string              `gorm:"column:mid"`
	Addition   string              `gorm:"column:addition"`
	Disable    bool                `gorm:"column:disable"`

	Acquirer *agency.Acquirer   `gorm:"-"`
	Terminal *acquirer.Terminal `gorm:"-"`
}

func (MerchantAccount) TableName() string {
	return "payment_merchant_account"
}

func (m *MerchantAccount) Get(merchantID uint64) (*MerchantAccount, error) {
	var ret = new(MerchantAccount)
	err := models.DB().Model(m).Where("id=?", merchantID).Find(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func (m *MerchantAccount) GetByAcquirerID(acqID, offset, limit uint64) ([]*MerchantAccount, error) {
	var ret []*MerchantAccount
	err := models.DB().Model(m).Where("acquirer_id=?", acqID).Offset(offset).Limit(limit).Find(ret).Error
	if err != nil {
		return nil, err
	}

	return ret, nil
}
