package tms

import (
	"tpayment/models"

	"github.com/labstack/echo"

	"github.com/jinzhu/gorm"
)

// merchantdevice model
type DeviceModel struct {
	models.BaseModel
	Name string `gorm:"column:name"` // 外键
}

func (DeviceModel) TableName() string {
	return "tms_model"
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

// 根据Model 获取设备信息
func GetModelByID(db *models.MyDB, ctx echo.Context, id uint) (*DeviceModel, error) {
	ret := new(DeviceModel)

	err := db.Model(&DeviceModel{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}
