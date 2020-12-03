package tms

import (
	"strconv"
	"tpayment/models"
	"tpayment/modules"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UploadFile struct {
	models.BaseModel
	AgencyId uint64 `gorm:"column:agency_id" json:"agency_id"`

	DeviceSn string `gorm:"column:device_sn" json:"device_sn"`
	FileName string `gorm:"column:file_name" json:"file_name"`
	FileUrl  string `gorm:"column:file_url" json:"file_url"`
}

func (UploadFile) TableName() string {
	return "tms_upload_file"
}

// 根据device ID获取设备信息
func GetUploadFileByID(db *models.MyDB, ctx *gin.Context, id uint64) (*UploadFile, error) {

	ret := new(UploadFile)

	err := db.Model(&UploadFile{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryUploadFileRecord(db *models.MyDB, ctx *gin.Context, offset, limit uint64, filters map[string]string) (uint64, []*UploadFile, error) {
	agency := modules.IsAgencyAdmin(ctx)

	equalData := make(map[string]string)
	if agency != nil { // 是机构管理员的话，就需要添加机构排查
		equalData["agency_id"] = strconv.FormatUint(uint64(agency.ID), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Model(&UploadFile{}).Where(sqlCondition)

	// 统计总数
	var total uint64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*UploadFile
	if err = tmpDb.Order("updated_at desc").Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
