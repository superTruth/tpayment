package tms

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"tpayment/models"
)

type DeviceInfo struct {
	gorm.Model
	AgencyId string `gorm:"column:agency_id"`

	DeviceSn    string `gorm:"column:device_sn"`
	DeviceCsn   string `gorm:"column:device_csn"`
	DeviceModel uint   `gorm:"column:device_model"`
	Alias       string `gorm:"column:alias"`

	RebootMode       int    `gorm:"column:reboot_mode"`
	RebootTime       string `gorm:"column:reboot_time"`
	RebootDayInWeek  int    `gorm:"column:reboot_day_in_week"`
	RebootDayInMonth int    `gorm:"column:reboot_day_in_month"`

	Power int `gorm:"column:power"`

	LocationLat string `gorm:"column:location_lat"`
	LocationLon string `gorm:"column:location_lon"`
	PushToken   string `gorm:"column:push_token"`

	CustomAttributes string `gorm:"column:custom_attributes"`
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

const (
	AppInDeviceExternalIdTypeDevice      = "merchantdevice"
	AppInDeviceExternalIdTypeBatchUpdate = "batch"
)

type AppInDevice struct {
	gorm.Model

	ExternalId     uint   `gorm:"column:external_id"`      // 外键
	ExternalIdType string `gorm:"column:external_id_type"` // 外键

	Name        string `gorm:"column:name"`
	PackageId   string `gorm:"column:package_id"`
	VersionName string `gorm:"column:version_name"`
	VersionCode int    `gorm:"column:version_code"`
	Status      string `gorm:"column:status"`

	AppID     uint `gorm:"column:app_id"`
	AppFileId uint `gorm:"column:app_file_id"`

	App     *App     `gorm:"-"`
	AppFile *AppFile `gorm:"-"`
}

func (AppInDevice) TableName() string {
	return "mdm2_app_in_device"
}

// 根据Device SN 获取设备信息
func GetDevice(deviceSn string) (*DeviceInfo, error) {

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

// 批量获取device
func GetDevices(deviceSns []string) ([]DeviceInfo, error) {
	var ret []DeviceInfo

	sb := strings.Builder{}
	for i, v := range deviceSns {
		sb.WriteString(v)
		if i == len(deviceSns)-1 { // 最后一个不用添加，
			break
		}
		sb.WriteString(",")
	}

	// 查找出所有device id
	rows, err := models.DB().Model(&_DeviceAndTagMid{}).
		Where("device_sn in ( ? )", sb.String()).
		Rows()

	defer rows.Close()

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return ret, nil
		}
		return nil, err
	}

	err = rows.Scan(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// 根据Device SN 获取app信息
func GetAppsInDevice(externalId uint, externalIdType string, offset int, limit int) ([]AppInDevice, error) {
	var ret []AppInDevice

	rows, err := models.DB().Raw("SELECT a.id, a.name, a.package_id, a.version_name, a.version_code, a.status, a.app_id, a.app_file_id, "+
		"b.name, b.package_id, "+
		"c.version_name, c.version_code, c.update_description, c.file_name, c.file_url "+
		"FROM mdm2_app_in_device a "+
		"LEFT JOIN mdm2_apps b ON a.app_id = b.id and b.deleted_at is null "+
		"LEFT JOIN mdm2_app_files c ON a.app_file_id = c.id and c.deleted_at is null "+
		"where a.external_id=? and a.external_id_type=? and a.deleted_at is null limit ? offset ?", externalId, externalIdType, limit, offset).Rows()
	if err != nil {
		fmt.Println("error->", err.Error())
		return ret, err
	}
	defer rows.Close()

	for rows.Next() {

		appInDevice := AppInDevice{
			App:     new(App),
			AppFile: new(AppFile),
		}

		err := rows.Scan(&appInDevice.ID, &appInDevice.Name, &appInDevice.PackageId, &appInDevice.VersionName, &appInDevice.VersionCode, &appInDevice.Status, &appInDevice.AppID, &appInDevice.AppFileId,
			&appInDevice.App.Name, &appInDevice.App.PackageId,
			&appInDevice.AppFile.VersionName,
			&appInDevice.AppFile.VersionCode, &appInDevice.AppFile.UpdateDescription, &appInDevice.AppFile.FileName, &appInDevice.AppFile.FileUrl)

		if err != nil {
			fmt.Println("err->", err)
		}

		// 只有这2种状态是有配置数据
		//if appInDevice.Status == nil{
		//	appInDevice.Status = utils.Int2PInt(conf.STATUS_PENDING_INSTALL)
		//}
		//if *appInDevice.Status != conf.STATUS_PENDING_INSTALL && *appInDevice.Status != conf.STATUS_INSTALLED && *appInDevice.Status != conf.STATUS_WARNING_INSTALLED {
		//	appInDevice.App = nil
		//	appInDevice.AppFile = nil
		//}

		if appInDevice.AppID == 0 || appInDevice.AppID == 0 {
			appInDevice.App = nil
		}
		if appInDevice.AppFileId == 0 || appInDevice.AppFileId == 0 {
			appInDevice.AppFile = nil
		}

		//// 这2种情况，如果没有实际的配置文件的话，就直接跳过
		//if (*appInDevice.Status == conf.STATUS_PENDING_INSTALL || *appInDevice.Status == conf.STATUS_INSTALLED || *appInDevice.Status == conf.STATUS_WARNING_INSTALLED) &&
		//	(appInDevice.App == nil || appInDevice.AppFile == nil) {
		//	fmt.Println("如果没有实际的配置文件的话")
		//	continue
		//}
		//data, _ := json.MarshalIndent(appInDevice, "", "   ")
		//fmt.Println("data->", string(data))

		ret = append(ret, appInDevice)
	}

	return ret, nil
}

// 创建设备
func CreateDevice(device *DeviceInfo) error {
	err := models.DB().Create(device).Error

	if err != nil {
		return err
	}

	return nil
}

// 删除设备
func DeleteDevice(deviceSn string) error {

	err := models.DB().Delete(&DeviceInfo{DeviceSn: deviceSn}).Error

	if err != nil {
		return err
	}

	return nil
}

// 更新设备信息
func UpdateDevice(device *DeviceInfo) error {
	err := models.DB().Model(device).Updates(device).Error

	if err != nil {
		return err
	}
	return nil
}

// 更新app in merchantdevice
func UpdateAppInDevice(appInDevice *AppInDevice) error {
	err := models.DB().Model(appInDevice).Updates(appInDevice).Error
	if err != nil {
		return err
	}

	return nil
}

// 删除掉app in merchantdevice
func DeleteAppInDevice(appInDevice *AppInDevice) error {
	err := models.DB().Delete(appInDevice).Error
	if err != nil {
		return err
	}

	return nil
}

func CreateAppInDevice(appInDevice *AppInDevice) error {
	err := models.DB().Create(appInDevice).Error
	if err != nil {
		return err
	}

	return nil
}

type App struct {
	gorm.Model

	Name        string `gorm:"column:name"`
	PackageId   string `gorm:"column:package_id"`
	Description string `gorm:"column:description"`
}

func (App) TableName() string {
	return "mdm2_apps"
}

type AppFile struct {
	gorm.Model

	VersionName       string   `gorm:"column:version_name"`
	VersionCode       int      `gorm:"column:version_code"`
	UpdateDescription string   `gorm:"column:update_description"`
	FileName          string   `gorm:"column:file_name"`
	FileUrl           []string `gorm:"column:file_url"`
	Status            int      `gorm:"column:decode_status"`
	DecodeFailMsg     string   `gorm:"column:decode_fail_msg"`

	AppId *uint `gorm:"column:app_id"`
}

func (AppFile) TableName() string {
	return "mdm2_app_files"
}

// 更新AppFile
func UpdateAppFile(appFile *AppFile) error {
	err := models.DB().Model(appFile).Updates(appFile).Error
	if err != nil {
		return err
	}

	return nil
}

func GetApkFileRecord(id uint) (*AppFile, error) {
	var ret = new(AppFile)

	err := models.DB().Where(&AppFile{Model: gorm.Model{ID: id}}).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}

		return nil, err
	}

	// TODO
	//if ret.Status == nil {
	//	ret.Status = 1
	//}
	//
	//if ret.VersionCode == nil {
	//	ret.VersionCode = utils.Int2PInt(0)
	//}
	//
	//if ret.VersionName == nil {
	//	ret.VersionName = utils.String2PString("")
	//}

	return ret, nil
}

// 获取app信息
func GetApp(id uint) (*App, error) {
	var ret = new(App)

	err := models.DB().Where(&App{Model: gorm.Model{ID: id}}).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}

		return nil, err
	}

	return ret, nil
}

// devicetag
type DeviceTag struct {
	gorm.Model

	AgencyId string  `gorm:"column:agency_id"`
	Name     *string `gorm:"column:name"` // 外键
}

type _DeviceAndTagMid struct {
	gorm.Model

	TagID    *uint `gorm:"column:tag_id"`
	DeviceId *uint `gorm:"column:device_id"`
}

func (_DeviceAndTagMid) TableName() string {
	return "mdm2_device_and_tag_mid"
}

func GetDeviceByTag(tagsUuid []string, offset int, limit int) ([]string, error) {

	var deviceUuids []string

	// 全选
	if len(tagsUuid) == 0 {
		err := models.DB().Model(&_DeviceAndTagMid{}).
			Offset(offset).Limit(offset).
			Select("device_sn").
			Find(deviceUuids).Error

		if err != nil {
			if gorm.ErrRecordNotFound == err { // 没有记录
				return deviceUuids, nil
			}
			return nil, err
		}

	}

	// 部分选择
	sb := strings.Builder{}
	for i, v := range tagsUuid {
		sb.WriteString(v)
		if i == len(tagsUuid)-1 { // 最后一个不用添加，
			break
		}
		sb.WriteString(",")
	}

	// 查找出所有device id
	rows, err := models.DB().Model(&_DeviceAndTagMid{}).
		Where("tag_uuid in ( ? )", sb.String()).
		Select("device_sn").
		Offset(offset).Limit(limit).
		Rows()

	defer rows.Close()

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return deviceUuids, nil
		}
		return nil, err
	}

	err = rows.Scan(&deviceUuids)
	if err != nil {
		return nil, err
	}

	return deviceUuids, nil
}

// merchantdevice model
type DeviceModel struct {
	gorm.Model

	Name *string `gorm:"column:name"` // 外键
}

func (DeviceModel) TableName() string {
	return "mdm2_models"
}

func GetModels() ([]DeviceModel, error) {
	var deviceModels []DeviceModel

	if err := models.DB().Model(&DeviceModel{}).Find(&deviceModels).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return deviceModels, nil
		}
		return deviceModels, err
	}

	return deviceModels, nil
}
