package tms

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"tpayment/models"
)

// 根据device ID获取设备信息
func GetAppInDeviceByID(db *models.MyDB, ctx echo.Context, id uint) (*AppInDevice, error) {

	ret := new(AppInDevice)

	err := db.Model(&AppInDevice{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryAppInDeviceRecord(db *models.MyDB, ctx echo.Context, deviceId ,offset, limit uint, filters map[string]string) (uint, []AppInDevice, error) {
	filterTmp := make(map[string]interface{})

	for k, v := range filters {
		filterTmp[k] = v
	}

	filterTmp["external_id"] = deviceId
	filterTmp["external_id_type"] = AppInDeviceExternalIdTypeDevice

	// conditions
	tmpDb := db.Table("mdm2_device_infos").Where(filterTmp)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []AppInDevice
	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
