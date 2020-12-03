package key

import (
	"tpayment/models"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
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

func (key *ApplePayKey) GetKeyByHash(db *models.MyDB, ctx *gin.Context, publicHash string) (*ApplePayKey, error) {
	ret := new(ApplePayKey)
	err := db.Model(key).Where("public_key_hash=?", publicHash).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func (key *ApplePayKey) GetKeyByDomain(db *models.MyDB, ctx *gin.Context, domain string) (*ApplePayKey, error) {
	ret := new(ApplePayKey)
	err := db.Model(key).Where("domain=?", domain).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}
