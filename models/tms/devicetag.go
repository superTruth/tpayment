package tms

import (
	"errors"
	"strconv"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/modules"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DeviceTagDao = &DeviceTag{}

type DeviceTag struct {
	models.BaseModel

	AgencyId    uint64 `json:"agency_id" gorm:"column:agency_id"`
	Name        string `json:"name" gorm:"column:name"` // 外键
	Description string `json:"description" gorm:"column:description"`
}

func (DeviceTag) TableName() string {
	return "tms_tags"
}

// 根据device ID获取设备信息
func GetDeviceTagByID(id uint64) (*DeviceTag, error) {

	ret := new(DeviceTag)

	err := models.DB.Model(&DeviceTag{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetDeviceTagByIDs(ids *models.IntArray) ([]*DeviceTag, error) {

	var ret []*DeviceTag

	err := models.DB.Model(&DeviceTag{}).Where("id IN (?)", ids.Change2UintArray()).Find(&ret).Error

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func QueryDeviceTagRecord(ctx *gin.Context, offset, limit uint64, filters map[string]string) (uint64, []*DeviceTag, error) {

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
	tmpDb := models.DB.Model(&DeviceTag{}).Where(sqlCondition)

	// 统计总数
	var total int64 = 0
	err = tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*DeviceTag
	if err = tmpDb.Order("id desc").Offset(int(offset)).Limit(int(limit)).Find(&ret).Error; err != nil {
		return uint64(total), ret, err
	}

	return uint64(total), ret, nil
}

func IsTagUsing(tagId uint64) (bool, error) {
	ret := new(DeviceAndTagMid)
	err := models.DB.Model(&DeviceAndTagMid{}).Where("tag_id=?", tagId).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (d *DeviceTag) GetInAgency(agencyID uint64) ([]*DeviceTag, error) {
	var ret []*DeviceTag

	err := models.DB.Model(d).Where("agency_id = ?", agencyID).Find(&ret).Error

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (d *DeviceTag) Create(tag *DeviceTag) error {
	return models.DB.Create(tag).Error
}

func (d *DeviceTag) QueryDeviceInTag(tagID, offset, limit uint64, filters map[string]string) (uint64, []*DeviceInfo, error) {
	var ret []*DeviceInfo

	equalData := make(map[string]string)
	sqlCondition := models.CombQueryCondition(equalData, filters)

	tmpDb := models.DB.Model(&DeviceInfo{}).Where(sqlCondition).
		Joins("join tms_device_and_tag_mid mid on tms_device.id = mid.device_id and mid.tag_id =? and mid.deleted_at is null", tagID)

	// 统计总数
	var total int64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	//
	err = tmpDb.Select("tms_device.*").Offset(int(offset)).Limit(int(limit)).Order("id desc").Find(&ret).Error
	if err != nil {
		return 0, nil, err
	}

	return uint64(total), ret, nil
}
