package tms

import (
	"github.com/jinzhu/gorm"
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
