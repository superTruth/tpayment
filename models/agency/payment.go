package agency

import (
	"tpayment/models"

	"gorm.io/gorm"
)

type PaymentMethod struct {
	gorm.Model

	Name string `json:"name" gorm:"column:name"`
}

func (PaymentMethod) TableName() string {
	return "agency_payment_method"
}

type PaymentType struct {
	gorm.Model

	Name string `json:"name" gorm:"column:name"`
}

func (PaymentType) TableName() string {
	return "agency_payment_type"
}

type EntryType struct {
	gorm.Model

	Name string `json:"name" gorm:"column:name"`
}

func (EntryType) TableName() string {
	return "agency_entry_type"
}

func GetPaymentMethods() ([]string, error) {
	var ret []PaymentMethod
	err := models.DB.Model(&PaymentMethod{}).Find(&ret).Error

	var retStr []string

	for _, v := range ret {
		retStr = append(retStr, v.Name)
	}

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return retStr, nil
}

func GetPaymentTypes() ([]string, error) {
	var ret []PaymentType
	err := models.DB.Model(&PaymentType{}).Find(&ret).Error

	var retStr []string

	for _, v := range ret {
		retStr = append(retStr, v.Name)
	}

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return retStr, nil
}

func GetEntryTypes() ([]string, error) {
	var ret []EntryType
	err := models.DB.Model(&EntryType{}).Find(&ret).Error

	var retStr []string

	for _, v := range ret {
		retStr = append(retStr, v.Name)
	}

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return retStr, nil
}
