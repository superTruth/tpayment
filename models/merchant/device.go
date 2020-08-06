package merchant

import (
	"tpayment/models"
	"tpayment/models/tms"

	"github.com/jinzhu/gorm"

	"github.com/labstack/echo"
)

type DeviceInMerchant struct {
	models.BaseModel

	DeviceId   uint `json:"device_id" gorm:"column:device_id"`
	MerchantId uint `json:"merchant_id" gorm:"column:merchant_id"`
}

func (DeviceInMerchant) TableName() string {
	return "device_in_merchant"
}

func GetDeviceInMerchantAssociateById(db *models.MyDB, ctx echo.Context, id uint) (*DeviceInMerchant, error) {
	ret := new(DeviceInMerchant)

	err := db.Model(&DeviceInMerchant{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

type DeviceInMerchantQueryBean struct {
	models.BaseModel
	DeviceId uint   `json:"device_id"`
	DeviceSn string `json:"device_sn"`
}

func QueryMerchantDeviceRecord(db *models.MyDB, ctx echo.Context, merchantId, offset, limit uint, filters map[string]string) (uint, []*DeviceInMerchantQueryBean, error) {
	// conditions
	tmpDb := db.Table(tms.DeviceInfo{}.TableName()).Model(&tms.DeviceInfo{})
	tmpDb = tmpDb.Joins("JOIN device_in_merchant ass ON ass.device_id = tms_device.id AND ass.merchant_id = ? AND ass.deleted_at IS NULL", merchantId)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*DeviceInMerchantQueryBean
	if err = tmpDb.Offset(offset).Limit(limit).
		Select("ass.id as id, ass.created_at as created_at, ass.updated_at as updated_at, tms_device.id as device_id, tms_device.device_sn as device_sn").
		Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
