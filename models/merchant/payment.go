package merchant

import "github.com/jinzhu/gorm"

type PaymentSettingInDevice struct {
	gorm.Model

	MerchantDeviceId uint     `json:"merchant_device_id" gorm:"column:merchant_device_id"`
	PaymentMethods   []string `json:"payment_methods" gorm:"column:payment_methods"`
	EntryTypes       []string `json:"entry_types" gorm:"column:entry_types"`
	PaymentTypes     []string `json:"payment_types" gorm:"column:payment_types"`
	AcquirerId       uint     `json:"acquirer_id" gorm:"column:acquirer_id"`
	Mid              string   `json:"mid" gorm:"column:mid"`
	Tid              string   `json:"tid" gorm:"column:tid"`
	Addition         string   `json:"addition" gorm:"column:addition"`
}

func (PaymentSettingInDevice) TableName() string {
	return "payment_setting_in_device"
}


