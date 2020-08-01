package tms

import (
	"errors"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/models/agency"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type DeviceInfo struct {
	gorm.Model
	AgencyId string `gorm:"column:agency_id" json:"agency_id"`

	DeviceSn    string `gorm:"column:device_sn" json:"device_sn"`
	DeviceCsn   string `gorm:"column:device_csn" json:"device_csn"`
	DeviceModel uint   `gorm:"column:device_model" json:"device_model"`
	Alias       string `gorm:"column:alias" json:"alias"`

	RebootMode       int    `gorm:"column:reboot_mode" json:"reboot_mode"`
	RebootTime       string `gorm:"column:reboot_time" json:"reboot_time"`
	RebootDayInWeek  int    `gorm:"column:reboot_day_in_week" json:"reboot_day_in_week"`
	RebootDayInMonth int    `gorm:"column:reboot_day_in_month" json:"reboot_day_in_month"`

	Battery int `gorm:"column:battery" json:"battery"`

	LocationLat string `gorm:"column:location_lat" json:"location_lat"`
	LocationLon string `gorm:"column:location_lon" json:"location_lon"`
	PushToken   string `gorm:"column:push_token" json:"push_token"`

	Tags []DeviceTagFull `gorm:"column:-" json:"tags"`
}

func (DeviceInfo) TableName() string {
	return "mdm2_device_infos"
}

func GenerateDeviceInfo() *DeviceInfo {
	device := new(DeviceInfo)

	device.RebootMode = 1
	device.RebootTime = "03:00"

	return device
}

// 根据Device SN 获取设备信息
func GetDeviceBySn(deviceSn string) (*DeviceInfo, error) {

	deviceInfo := new(DeviceInfo)

	err := models.DB().Where(&DeviceInfo{DeviceSn: deviceSn}).First(&deviceInfo).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}

		return nil, err
	}

	return deviceInfo, nil
}

// 根据device ID获取设备信息
func GetDeviceByID(db *models.MyDB, ctx echo.Context, id uint) (*DeviceInfo, error) {

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

func QueryDeviceRecord(db *models.MyDB, ctx echo.Context, offset, limit uint, filters map[string]string) (uint, []DeviceInfo, error) {
	filterTmp := make(map[string]interface{})
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
	agencys := ctx.Get(conf.ContextTagAgency).([]*agency.Agency)

	for k, v := range filters {
		filterTmp[k] = v
	}

	if userBean.Role != string(conf.RoleAdmin) { // 管理员，不需要过滤机构
		if len(agencys) == 0 {
			return 0, nil, errors.New("user not agency admin")
		}
		filterTmp["agency_id"] = agencys[0].ID
	}

	// conditions
	tmpDb := db.Model(&DeviceInfo{}).Where(filterTmp)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []DeviceInfo
	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

func QueryTags(db *models.MyDB, ctx echo.Context, device *DeviceInfo) ([]DeviceTagFull, error) {
	var ret []DeviceTagFull
	filterTmp := make(map[string]interface{})
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
	agencys := ctx.Get(conf.ContextTagAgency).([]*agency.Agency)
	if userBean.Role != string(conf.RoleAdmin) { // 管理员，不需要过滤机构
		if len(agencys) == 0 {
			return ret, errors.New("user not agency admin")
		}
		filterTmp["agency_id"] = agencys[0].ID
	}

	err := db.Model(&DeviceTag{}).Joins("JOIN mdm2_device_and_tag_mid mid ON mid.device_id=? AND mid.tag_id=mdm2_tags.id and mid.deleted_at is null", device.ID).
		Where(filterTmp).
		Select("mdm2_tags.id as id, mdm2_tags.agency_id as agency_id, mdm2_tags.name as name, mdm2_tags.created_at as created_at, mdm2_tags.updated_at as updated_at, mid.id as mid_id").
		Find(&ret).Error

	if err != nil {
		return ret, err
	}

	return ret, nil
}
