package merchant

import "github.com/jinzhu/gorm"

type DeviceInMerchant struct {
	gorm.Model

	DeviceId   uint `json:"device_id" gorm:"column:device_id"`
	MerchantId uint `json:"merchant_id" gorm:"column:merchant_id"`
}

func (DeviceInMerchant) TableName() string {
	return "device_in_merchant"
}

type PaymentSettingInDevice struct {
	gorm.Model

	DeviceInMerchantId uint `json:"device_in_merchant_id" gorm:"column:device_in_merchant_id"`

	PaymentMethods []string `json:"payment_methods" gorm:"column:payment_methods"`
	EntryTypes     []string `json:"entry_types" gorm:"column:entry_types"`
	PaymentTypes   []string `json:"payment_types" gorm:"column:payment_types"`

	AcquirerTypeId uint   `json:"acquirer_type_id" gorm:"column:acquirer_type_id"`
	MID            string `json:"mid" gorm:"column:mid"`
	TID            string `json:"tid" gorm:"column:tid"`
	Addition       string `json:"addition" gorm:"column:addition"`
}

func (PaymentSettingInDevice) TableName() string {
	return "payment_setting_in_device"
}
