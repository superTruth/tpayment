package merchant

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"tpayment/models"
)

type DeviceInMerchant struct {
	gorm.Model

	DeviceId   uint `json:"device_id" gorm:"column:device_id"`
	MerchantId uint `json:"merchant_id" gorm:"column:merchant_id"`
}

func (DeviceInMerchant) TableName() string {
	return "device_in_merchant"
}

type DeviceInMerchantQueryBean struct {
	gorm.Model
	DeviceId uint   `json:"device_id"`
	DeviceSn string `json:"device_sn"`
}

func QueryMerchantDeviceRecord(db *models.MyDB, ctx echo.Context, merchantId, offset, limit uint, filters map[string]string) (uint, []DeviceInMerchantQueryBean, error) {
	filterTmp := make(map[string]interface{})

	for k, v := range filters {
		filterTmp[k] = v
	}

	if merchantId != 0 {
		filterTmp["merchant_id"] = merchantId
	}

	// conditions
	tmpDb := db.Table("device").Where(filterTmp)
	tmpDb = tmpDb.Joins("JOIN device_in_merchant ass ON ass.device_id = device.id AND ass.merchant_id = ? AND ass.deleted_at IS NULL", merchantId)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []DeviceInMerchantQueryBean
	if err = tmpDb.Offset(offset).Limit(limit).
		Select("ass.id as id, ass.created_at as created_at, ass.updated_at as updated_at, device.id as device_id, device.sn as device_sn").
		Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
