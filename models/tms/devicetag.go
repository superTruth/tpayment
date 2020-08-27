package tms

import (
	"errors"
	"strconv"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/modules"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type DeviceTag struct {
	models.BaseModel

	AgencyId uint   `json:"agency_id" gorm:"column:agency_id"`
	Name     string `json:"name" gorm:"column:name"` // 外键
}

func (DeviceTag) TableName() string {
	return "tms_tags"
}

// 根据device ID获取设备信息
func GetDeviceTagByID(db *models.MyDB, ctx *gin.Context, id uint) (*DeviceTag, error) {

	ret := new(DeviceTag)

	err := db.Model(&DeviceTag{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetDeviceTagByIDs(db *models.MyDB, ctx *gin.Context, ids *models.IntArray) ([]*DeviceTag, error) {

	var ret []*DeviceTag

	err := db.Model(&DeviceTag{}).Where("id IN (?)", ids.Change2UintArray()).Find(&ret).Error

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func QueryDeviceTagRecord(db *models.MyDB, ctx *gin.Context, offset, limit uint, filters map[string]string) (uint, []*DeviceTag, error) {

	agencyId, err := modules.GetAgencyId2(ctx)
	if err != nil {
		return 0, nil, errors.New(conf.NoPermission.String())
	}
	equalData := make(map[string]string)
	if agencyId != 0 {
		equalData["agency_id"] = strconv.FormatUint(uint64(agencyId), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Model(&DeviceTag{}).Where(sqlCondition)

	// 统计总数
	var total uint = 0
	err = tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*DeviceTag
	if err = tmpDb.Order("updated_at desc").Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
