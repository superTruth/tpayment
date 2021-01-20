package tms

import (
	"errors"
	"strconv"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/modules"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var DeviceInfoDao = &DeviceInfo{}

type DeviceInfo struct {
	models.BaseModel
	AgencyId uint64 `gorm:"column:agency_id" json:"agency_id"`

	DeviceSn        string `gorm:"column:device_sn" json:"device_sn"`
	DeviceCsn       string `gorm:"column:device_csn" json:"device_csn"`
	DeviceModel     uint64 `gorm:"column:device_model" json:"-"`
	DeviceModelName string `gorm:"-" json:"device_model"`
	Alias           string `gorm:"column:alias" json:"alias"`

	RebootMode       string `gorm:"column:reboot_mode" json:"reboot_mode"`
	RebootTime       string `gorm:"column:reboot_time" json:"reboot_time"`
	RebootDayInWeek  int    `gorm:"column:reboot_day_in_week" json:"reboot_day_in_week"`
	RebootDayInMonth int    `gorm:"column:reboot_day_in_month" json:"reboot_day_in_month"`

	Battery int `gorm:"column:battery" json:"battery"`

	LocationLat string `gorm:"column:location_lat" json:"location_lat"`
	LocationLon string `gorm:"column:location_lon" json:"location_lon"`
	PushToken   string `gorm:"column:push_token" json:"push_token"`

	Tags *[]*DeviceTagFull `gorm:"-" json:"tags,omitempty"`
}

func (DeviceInfo) TableName() string {
	return "tms_device"
}

func GenerateDeviceInfo() *DeviceInfo {
	device := new(DeviceInfo)

	device.RebootMode = conf.RebootModeEveryDay
	device.RebootTime = "03:00"

	return device
}

// 根据Device SN 获取设备信息
func GetDeviceBySn(db *models.MyDB, ctx *gin.Context, deviceSn string) (*DeviceInfo, error) {

	deviceInfo := new(DeviceInfo)

	err := db.Where(&DeviceInfo{DeviceSn: deviceSn}).First(&deviceInfo).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}

		return nil, err
	}

	return deviceInfo, nil
}

func (d *DeviceInfo) GetBySn(deviceSn string) (*DeviceInfo, error) {
	deviceInfo := new(DeviceInfo)

	err := models.DB().Model(d).Where(&DeviceInfo{DeviceSn: deviceSn}).First(&deviceInfo).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}

		return nil, err
	}

	return deviceInfo, nil
}

// 根据device ID获取设备信息
func GetDeviceByID(db *models.MyDB, ctx *gin.Context, id uint64) (*DeviceInfo, error) {

	ret := new(DeviceInfo)

	err := db.Model(&DeviceInfo{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func ResetDeviceAgency(db *models.MyDB, ctx *gin.Context, id uint64) error {
	err := db.Model(&DeviceInfo{}).Where("id=?", id).Update("agency_id", 0).Error

	if err != nil {
		return err
	}

	return nil
}

func QueryDeviceRecordByAgencyId(db *models.MyDB, ctx *gin.Context, agencyId, offset, limit uint64, filters map[string]string) (uint64, []*DeviceInfo, error) {
	equalData := make(map[string]string)
	equalData["agency_id"] = strconv.FormatUint(uint64(agencyId), 10)
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Model(&DeviceInfo{}).Where(sqlCondition)

	// 统计总数
	var total uint64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*DeviceInfo
	if err = tmpDb.Order("updated_at desc").Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

func QueryDeviceRecord(db *models.MyDB, ctx *gin.Context, offset, limit uint64, filters map[string]string) (uint64, []*DeviceInfo, error) {

	agency := modules.IsAgencyAdmin(ctx)

	equalData := make(map[string]string)
	if agency != nil { // 是机构管理员的话，就需要添加机构排查
		equalData["agency_id"] = strconv.FormatUint(uint64(agency.ID), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Model(&DeviceInfo{}).Where(sqlCondition)

	// 统计总数
	var total uint64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*DeviceInfo
	if err = tmpDb.Order("updated_at desc").Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

func QueryTagsInDevice(db *models.MyDB, ctx *gin.Context, device *DeviceInfo) ([]*DeviceTagFull, error) {
	var ret []*DeviceTagFull
	filterTmp := make(map[string]interface{})
	if modules.IsAdmin(ctx) == nil {
		agency := modules.IsAgencyAdmin(ctx)
		if agency == nil {
			return nil, errors.New(conf.NoPermission.String())
		}
		filterTmp["agency_id"] = agency.ID
	}

	err := db.Table(DeviceTag{}.TableName()).Model(&DeviceTag{}).Joins("JOIN tms_device_and_tag_mid mid ON mid.device_id=? AND mid.tag_id=tms_tags.id and mid.deleted_at is null", device.ID).
		Where(filterTmp).Order("updated_at desc").
		Select("tms_tags.id as id, tms_tags.agency_id as agency_id, tms_tags.name as name, tms_tags.created_at as created_at, tms_tags.updated_at as updated_at, mid.id as mid_id").
		Find(&ret).Error

	if err != nil {
		return ret, err
	}

	return ret, nil
}

// devicetag
type DeviceTagFull struct {
	DeviceTag
	MidId uint64 `json:"agency_id" gorm:"column:mid_id"`
}

type DeviceAndTagMid struct {
	models.BaseModel

	TagID    uint64 `gorm:"column:tag_id"`
	DeviceId uint64 `gorm:"column:device_id"`
}

func (DeviceAndTagMid) TableName() string {
	return "tms_device_and_tag_mid"
}
