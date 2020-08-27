package tms

import (
	"tpayment/models"

	"github.com/gin-gonic/gin"
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
func GetModelByID(db *models.MyDB, ctx *gin.Context, id uint) (*DeviceModel, error) {
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

func GetModelByIDs(db *models.MyDB, ctx *gin.Context, ids *models.IntArray) ([]*DeviceModel, error) {
	var ret []*DeviceModel

	err := db.Model(&DeviceModel{}).Where("id IN (?)", ids.Change2UintArray()).Find(&ret).Error

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func QueryModelRecord(db *models.MyDB, ctx *gin.Context, offset, limit uint, filters map[string]string) (uint, []*DeviceModel, error) {
	equalData := make(map[string]string)
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Model(&DeviceModel{}).Where(sqlCondition)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*DeviceModel
	if err = tmpDb.Order("updated_at desc").Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
