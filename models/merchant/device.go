package merchant

import (
	"tpayment/models"
	"tpayment/models/tms"

	"gorm.io/gorm"
)

var DeviceInMerchantDao = &DeviceInMerchant{}

type DeviceInMerchant struct {
	models.BaseModel

	DeviceId   uint64 `json:"device_id" gorm:"column:device_id"`
	MerchantId uint64 `json:"merchant_id" gorm:"column:merchant_id"`
}

func (DeviceInMerchant) TableName() string {
	return "merchant_device"
}

func (d *DeviceInMerchant) GetByMerchantIdAndDeviceID(merchantID, deviceID uint64) (*DeviceInMerchant, error) {
	ret := new(DeviceInMerchant)

	err := models.DB.Model(d).Where("device_id=? and merchant_id=?", deviceID, merchantID).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetDeviceInMerchantAssociateById(id uint64) (*DeviceInMerchant, error) {
	ret := new(DeviceInMerchant)

	err := models.DB.Model(&DeviceInMerchant{}).Where("id=?", id).First(ret).Error

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
	DeviceId  uint64 `json:"device_id"`
	DeviceSn  string `json:"device_sn"`
	DeviceCsn string `json:"device_csn"`
}

func QueryMerchantDeviceRecord(merchantId, offset, limit uint64, filters map[string]string) (uint64, []*DeviceInMerchantQueryBean, error) {
	// conditions
	tmpDb := models.DB.Table(tms.DeviceInfo{}.TableName()).Model(&tms.DeviceInfo{})
	tmpDb = tmpDb.Joins("JOIN merchant_device ass ON ass.device_id = tms_device.id AND ass.merchant_id = ? AND ass.deleted_at IS NULL", merchantId)
	if _, ok := filters["device_sn"]; ok {
		tmpDb = tmpDb.Where("tms_device.device_sn LIKE ?", filters["device_sn"]+"%")
	}

	// 统计总数
	var total int64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*DeviceInMerchantQueryBean
	if err = tmpDb.Offset(int(offset)).Limit(int(limit)).
		Select("ass.id as id, ass.created_at as created_at, ass.updated_at as updated_at, tms_device.id as device_id, tms_device.device_sn as device_sn, tms_device.device_csn as device_csn").
		Find(&ret).Error; err != nil {
		return uint64(total), ret, err
	}

	return uint64(total), ret, nil
}
