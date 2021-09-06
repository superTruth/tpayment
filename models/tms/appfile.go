package tms

import (
	"strconv"
	"tpayment/models"

	"gorm.io/gorm"
)

type AppFile struct {
	models.BaseModel

	VersionName       string `gorm:"column:version_name" json:"version_name"`
	VersionCode       int    `gorm:"column:version_code" json:"version_code"`
	UpdateDescription string `gorm:"column:update_description" json:"update_description"`
	FileName          string `gorm:"column:file_name" json:"file_name"`
	FileUrl           string `gorm:"column:file_url" json:"file_url"`
	Status            string `gorm:"column:decode_status" json:"status"`
	DecodeFailMsg     string `gorm:"column:decode_fail_msg" json:"decode_fail_msg"`

	AppId uint64 `gorm:"column:app_id" json:"app_id"`
}

func (AppFile) TableName() string {
	return "tms_app_file"
}

// 根据device ID获取设备信息
func GetAppFileByID(id uint64) (*AppFile, error) {

	ret := new(AppFile)

	err := models.DB.Model(&AppFile{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryAppFileRecord(appId, offset, limit uint64, filters map[string]string) (uint64, []*AppFile, error) {
	equalData := make(map[string]string)
	equalData["app_id"] = strconv.FormatUint(uint64(appId), 10)
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := models.DB.Model(&AppFile{}).Where(sqlCondition)

	// 统计总数
	var total int64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*AppFile
	if err = tmpDb.Order("id desc").Offset(int(offset)).Limit(int(limit)).Find(&ret).Error; err != nil {
		return uint64(total), ret, err
	}

	return uint64(total), ret, nil
}
