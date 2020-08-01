package tms

import (
	"tpayment/models"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type AppInDevice struct {
	gorm.Model

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
	return "mdm2_app_in_device"
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
		return ret, err
	}
	// nolint
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
			return ret, err
		}

		// 只有这2种状态是有配置数据
		//if appInDevice.Status == nil{
		//	appInDevice.Status = utils.Int2PInt(conf.STATUS_PENDING_INSTALL)
		//}
		//if *appInDevice.Status != conf.STATUS_PENDING_INSTALL && *appInDevice.Status != conf.STATUS_INSTALLED && *appInDevice.Status != conf.STATUS_WARNING_INSTALLED {
		//	appInDevice.App = nil
		//	appInDevice.AppFile = nil
		//}

		if appInDevice.AppID == 0 {
			appInDevice.App = nil
		}
		if appInDevice.AppFileId == 0 {
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

func QueryAppInDeviceRecord(db *models.MyDB, ctx echo.Context, deviceId, offset, limit uint,
	externalType string, filters map[string]string) (uint, []AppInDevice, error) {
	filterTmp := make(map[string]interface{})

	for k, v := range filters {
		filterTmp[k] = v
	}

	filterTmp["external_id"] = deviceId
	filterTmp["external_id_type"] = externalType //AppInDeviceExternalIdTypeDevice

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
