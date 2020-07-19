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

func QueryAppInDeviceRecord(db *models.MyDB, ctx echo.Context, deviceId, offset, limit uint, filters map[string]string) (uint, []AppInDevice, error) {
	filterTmp := make(map[string]interface{})

	for k, v := range filters {
		filterTmp[k] = v
	}

	filterTmp["external_id"] = deviceId
	filterTmp["external_id_type"] = AppInDeviceExternalIdTypeDevice

	// conditions
	tmpDb := db.Model(&AppInDevice{}).Where(filterTmp)

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

	for i := 0; i < len(ret); i++ {
		// 没有配置，就不需要搜索
		if (ret[i].AppFileId == 0) || (ret[i].AppID == 0) {
			continue
		}

		// 获取app数据
		if ret[i].AppID != 0 {
			ret[i].App, err = GetAppByID(db, ctx, ret[i].AppID)
			if err != nil {
				return 0, nil, err
			}
		}
		// 获取app file数据
		if ret[i].AppFileId != 0 {
			ret[i].AppFile, err = GetAppFileByID(db, ctx, ret[i].AppFileId)
			if err != nil {
				return 0, nil, err
			}
		}
	}

	return total, ret, nil
}

// 1. 只有package id这种非法安装的app
// 2. 包含app id这种配置安装的app
func FindAppInDevice(db *models.MyDB, ctx echo.Context, appInDevice *AppInDevice) (*AppInDevice, error) {
	ret := new(AppInDevice)
	err := db.Model(&AppInDevice{}).Where("external_id=? AND external_id_type=? AND ((package_id=app.package_id) OR (app_id=app.id)) ",
		appInDevice.ExternalId, AppInDeviceExternalIdTypeDevice).Joins("mdm2_apps app ON app.id=? AND deleted_at is null", appInDevice.AppID).
		First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}
