package agency

import (
	"tpayment/models"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
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

func GetPaymentMethods(db *models.MyDB, ctx echo.Context) ([]string, error) {
	var ret []PaymentMethod
	err := db.Model(&PaymentMethod{}).Find(&ret).Error

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

func GetPaymentTypes(db *models.MyDB, ctx echo.Context) ([]string, error) {
	var ret []PaymentType
	err := db.Model(&PaymentType{}).Find(&ret).Error

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

func GetEntryTypes(db *models.MyDB, ctx echo.Context) ([]string, error) {
	var ret []EntryType
	err := db.Model(&EntryType{}).Find(&ret).Error

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
