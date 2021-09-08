package tms

import (
	"strconv"
	"tpayment/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppInDevice struct {
	models.BaseModel

	ExternalId     uint64 `gorm:"column:external_id" json:"external_id"`           // 外键
	ExternalIdType string `gorm:"column:external_id_type" json:"external_id_type"` // 外键

	Name        string `gorm:"column:name" json:"name"`
	PackageId   string `gorm:"column:package_id" json:"package_id"`
	VersionName string `gorm:"column:version_name" json:"version_name"`
	VersionCode int    `gorm:"column:version_code" json:"version_code"`
	Status      string `gorm:"column:status" json:"status"`

	AppID     uint64 `gorm:"column:app_id" json:"app_id"`
	AppFileId uint64 `gorm:"column:app_file_id" json:"app_file_id"`

	App     *App     `gorm:"-" json:"app"`
	AppFile *AppFile `gorm:"-" json:"app_file"`
}

func (AppInDevice) TableName() string {
	return "tms_app_in_device"
}

const (
	AppInDeviceExternalIdTypeDevice      = "merchantdevice"
	AppInDeviceExternalIdTypeBatchUpdate = "batch"
)

func GetAppsInDevice(externalId uint64, externalIdType string, offset uint64, limit uint64) (uint64, []*AppInDevice, error) {
	var ret []*AppInDevice

	equalData := make(map[string]string)
	equalData["external_id"] = strconv.FormatUint(uint64(externalId), 10)
	equalData["external_id_type"] = externalIdType
	sqlCondition := models.CombQueryCondition(equalData, make(map[string]string))

	tmpDb := models.DB.Model(&AppInDevice{}).Where(sqlCondition)

	// 统计总数
	var total int64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	if err = tmpDb.Order("id desc").Offset(int(offset)).Limit(int(limit)).Find(&ret).Error; err != nil {
		return uint64(total), ret, err
	}

	for i := 0; i < len(ret); i++ {
		if ret[i].AppID != 0 {
			ret[i].App, err = GetAppByID(ret[i].AppID)
			if err != nil {
				return 0, ret, err
			}
		}
		if ret[i].AppFileId != 0 {
			ret[i].AppFile, err = GetAppFileByID(ret[i].AppFileId)
			if err != nil {
				return 0, ret, err
			}
		}
	}

	return uint64(total), ret, nil
}

// 根据device ID获取设备信息
func GetAppInDeviceByID(id uint64) (*AppInDevice, error) {

	ret := new(AppInDevice)

	err := models.DB.Model(&AppInDevice{}).Where("id=?", id).First(ret).Error

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
func FindAppInDevice(ctx *gin.Context, deviceId uint64, appInDevice *AppInDevice) (*AppInDevice, error) {
	ret := new(AppInDevice)
	//err := models.DB.Model(&AppInDevice{}).Where("external_id=? AND external_id_type=? AND ((package_id=app.package_id) OR (app_id=app.id)) ",
	//	appInDevice.ExternalId, AppInDeviceExternalIdTypeDevice).Joins("tms_app app ON app.id=? AND deleted_at is null", appInDevice.AppID).
	//	First(ret).Error
	//if err != nil {
	//	if gorm.ErrRecordNotFound == err { // 没有记录
	//		return nil, nil
	//	}
	//	return nil, err
	//}

	err := models.DB.Model(&AppInDevice{}).Where("external_id=? AND external_id_type=? AND ((tms_app_in_device.package_id=app.package_id) OR (tms_app_in_device.app_id=app.id)) ",
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
		ret.App, err = GetAppByID(ret.AppID)
		if err != nil {
			return ret, err
		}
	}
	if ret.AppFileId != 0 {
		ret.AppFile, err = GetAppFileByID(ret.AppFileId)
		if err != nil {
			return ret, err
		}
	}

	return ret, nil
}
