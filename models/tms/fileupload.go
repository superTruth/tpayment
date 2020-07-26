package tms

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"tpayment/models"
)

type FileUpload struct {
	gorm.Model
	StoreID *int `gorm:"column:store_id"`

	DeviceSn *string `gorm:"column:device_sn"`
	FileName *string `gorm:"column:file_name"`
	FileUrl  *string `gorm:"column:file_url"`
}

func (FileUpload) TableName() string {
	return "mdm2_file_upload"
}

// 创建一条记录
func CreateFileUpload(bean *FileUpload) error {
	err := models.DB().Create(bean).Error

	if err != nil {
		return err
	}

	return nil
}

// 根据device ID获取设备信息
func GetUploadFileByID(db *models.MyDB, ctx echo.Context, id uint) (*FileUpload, error) {

	ret := new(FileUpload)

	err := db.Model(&FileUpload{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryUploadFileRecord(db *models.MyDB, ctx echo.Context, offset, limit uint, filters map[string]string) (uint, []FileUpload, error) {
	filterTmp := make(map[string]interface{})

	for k, v := range filters {
		filterTmp[k] = v
	}

	// conditions
	tmpDb := db.Table("mdm2_file_upload").Where(filterTmp)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []FileUpload
	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
