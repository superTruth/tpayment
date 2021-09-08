package tms

import (
	"tpayment/models"

	"gorm.io/gorm"
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

	if err := models.DB.Model(&DeviceModel{}).Find(&deviceModels).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return deviceModels, nil
		}
		return deviceModels, err
	}

	return deviceModels, nil
}

// 根据Model 获取设备信息
func GetModelByID(id uint64) (*DeviceModel, error) {
	ret := new(DeviceModel)

	err := models.DB.Model(&DeviceModel{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetModelByIDs(ids *models.IntArray) ([]*DeviceModel, error) {
	var ret []*DeviceModel

	err := models.DB.Model(&DeviceModel{}).Where("id IN (?)", ids.Change2UintArray()).Find(&ret).Error

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func QueryModelRecord(offset, limit uint64, filters map[string]string) (uint64, []*DeviceModel, error) {
	equalData := make(map[string]string)
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := models.DB.Model(&DeviceModel{}).Where(sqlCondition)

	// 统计总数
	var total int64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*DeviceModel
	if err = tmpDb.Order("id desc").Offset(int(offset)).Limit(int(limit)).Find(&ret).Error; err != nil {
		return uint64(total), ret, err
	}

	return uint64(total), ret, nil
}

func GetModelByName(name string) (*DeviceModel, error) {
	ret := new(DeviceModel)

	err := models.DB.Model(&DeviceModel{}).Where("name=?", name).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}
