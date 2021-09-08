package key

import (
	"tpayment/models"

	"gorm.io/gorm"
)

type ApplePayKey struct {
	models.BaseModel

	AgencyID      uint64 `gorm:"column:agency_id"`
	PublicKeyHash string `gorm:"column:public_key_hash"`
	Domain        string `gorm:"column:domain"`

	PublicKey     string `gorm:"column:public_key"`
	PrivateKey    string `gorm:"column:private_key"`
	TlsPublicKey  string `gorm:"column:tls_public_key"`
	TlsPrivateKey string `gorm:"column:tls_private_key"`
}

func (ApplePayKey) TableName() string {
	return "payment_apple_pay_key"
}

func (key *ApplePayKey) GetKeyByHash(publicHash string) (*ApplePayKey, error) {
	ret := new(ApplePayKey)
	err := models.DB.Model(key).Where("public_key_hash=?", publicHash).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func (key *ApplePayKey) GetKeyByDomain(domain string) (*ApplePayKey, error) {
	ret := new(ApplePayKey)
	err := models.DB.Model(key).Where("domain=?", domain).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}
