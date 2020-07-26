package tms

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"tpayment/models"
)

type AppFile struct {
	gorm.Model

	VersionName       string `gorm:"column:version_name"`
	VersionCode       int    `gorm:"column:version_code"`
	UpdateDescription string `gorm:"column:update_description"`
	FileName          string `gorm:"column:file_name"`
	FileUrl           string `gorm:"column:file_url"`
	Status            string `gorm:"column:decode_status"`
	DecodeFailMsg     string `gorm:"column:decode_fail_msg"`

	AppId *uint `gorm:"column:app_id"`
}

func (AppFile) TableName() string {
	return "mdm2_app_files"
}

// 根据device ID获取设备信息
func GetAppFileByID(db *models.MyDB, ctx echo.Context, id uint) (*AppFile, error) {

	ret := new(AppFile)

	err := db.Model(&AppFile{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryAppFileRecord(db *models.MyDB, ctx echo.Context, appId, offset, limit uint, filters map[string]string) (uint, []AppFile, error) {
	filterTmp := make(map[string]interface{})

	for k, v := range filters {
		filterTmp[k] = v
	}

	filterTmp["app_id"] = appId

	// conditions
	tmpDb := db.Table("mdm2_app_files").Where(filterTmp)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []AppFile
	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
