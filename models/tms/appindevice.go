package tms

import (
	"strconv"
	"tpayment/models"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type AppInDevice struct {
	models.BaseModel

	ExternalId     uint   `gorm:"column:external_id" json:"external_id"`           // 外键
	ExternalIdType string `gorm:"column:external_id_type" json:"external_id_type"` // 外键

	Name        string `gorm:"column:name" json:"name"`
	PackageId   string `gorm:"column:package_id" json:"package_id"`
	VersionName string `gorm:"column:version_name" json:"version_name"`
	VersionCode int    `gorm:"column:version_code" json:"version_code"`
	Status      string `gorm:"column:status" json:"status"`

	AppID     uint `gorm:"column:app_id" json:"app_id"`
	AppFileId uint `gorm:"column:app_file_id" json:"app_file_id"`

	App     *App     `gorm:"-" json:"app"`
	AppFile *AppFile `gorm:"-" json:"app_file"`
}

func (AppInDevice) TableName() string {
	return "tms_app_in_device"
}

func GetAppsInDevice(db *models.MyDB, ctx echo.Context, externalId uint, externalIdType string, offset uint, limit uint) (uint, []*AppInDevice, error) {
	var ret []*AppInDevice

	equalData := make(map[string]string)
	equalData["external_id"] = strconv.FormatUint(uint64(externalId), 10)
	equalData["external_id_type"] = externalIdType
	sqlCondition := models.CombQueryCondition(equalData, make(map[string]string))

	tmpDb := db.Model(&AppInDevice{}).Where(sqlCondition)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	for i := 0; i < len(ret); i++ {
		if ret[i].AppID != 0 {
			ret[i].App, err = GetAppByID(db, ctx, ret[i].AppID)
			if err != nil {
				return 0, ret, err
			}
		}
		if ret[i].AppFileId != 0 {
			ret[i].AppFile, err = GetAppFileByID(db, ctx, ret[i].AppFileId)
			if err != nil {
				return 0, ret, err
			}
		}
	}

	return total, ret, nil
}

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

// 1. 只有package id这种非法安装的app
// 2. 包含app id这种配置安装的app
func FindAppInDevice(db *models.MyDB, ctx echo.Context, deviceId uint, appInDevice *AppInDevice) (*AppInDevice, error) {
	ret := new(AppInDevice)
	//err := db.Model(&AppInDevice{}).Where("external_id=? AND external_id_type=? AND ((package_id=app.package_id) OR (app_id=app.id)) ",
	//	appInDevice.ExternalId, AppInDeviceExternalIdTypeDevice).Joins("tms_app app ON app.id=? AND deleted_at is null", appInDevice.AppID).
	//	First(ret).Error
	//if err != nil {
	//	if gorm.ErrRecordNotFound == err { // 没有记录
	//		return nil, nil
	//	}
	//	return nil, err
	//}

	err := db.Model(&AppInDevice{}).Where("external_id=? AND external_id_type=? AND ((tms_app_in_device.package_id=app.package_id) OR (tms_app_in_device.app_id=app.id)) ",
		deviceId, AppInDeviceExternalIdTypeDevice).
		Joins("JOIN tms_app app ON app.id=? AND app.deleted_at is null", appInDevice.AppID).
		First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	if ret.AppID != 0 {
		ret.App, err = GetAppByID(db, ctx, ret.AppID)
		if err != nil {
			return ret, err
		}
	}
	if ret.AppFileId != 0 {
		ret.AppFile, err = GetAppFileByID(db, ctx, ret.AppFileId)
		if err != nil {
			return ret, err
		}
	}

	return ret, nil
}
