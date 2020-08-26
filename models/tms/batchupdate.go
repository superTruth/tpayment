package tms

import (
	"strconv"
	"tpayment/models"
	"tpayment/modules"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type BatchUpdate struct {
	models.BaseModel

	AgencyId uint `gorm:"column:agency_id" json:"agency_id"`

	Description string `gorm:"column:description" json:"description"`
	Status      string `gorm:"column:status" json:"status"`

	UpdateFailMsg string `gorm:"column:update_fail_msg" json:"update_fail_msg"`

	Tags         *models.StringArray `gorm:"column:tags" json:"-"`
	DeviceModels *models.StringArray `gorm:"column:device_models" json:"device_models"`

	ConfigTags []*DeviceTag `gorm:"column:-" json:"tags"`

	Apps []*AppInDevice `gorm:"-" json:"-"`
}

func (BatchUpdate) TableName() string {
	return "tms_batch_update"
}

func GetBatchUpdateRecordById(db *models.MyDB, ctx *gin.Context, id uint) (*BatchUpdate, error) {
	ret := new(BatchUpdate)

	err := db.Model(&BatchUpdate{}).Where("id=?", id).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}

		return nil, err
	}

	return ret, nil
}

func QueryBatchUpdateRecord(db *models.MyDB, ctx *gin.Context, offset, limit uint, filters map[string]string) (uint, []*BatchUpdate, error) {
	agency := modules.IsAgencyAdmin(ctx)

	equalData := make(map[string]string)
	if agency != nil { // 是机构管理员的话，就需要添加机构排查
		equalData["agency_id"] = strconv.FormatUint(uint64(agency.ID), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Model(&BatchUpdate{}).Where(sqlCondition)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*BatchUpdate
	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

func GetBatchUpdateDevices(db *models.MyDB, ctx *gin.Context, batchUpdate *BatchUpdate, offset uint, limit uint) ([]*DeviceInfo, error) {
	tmpDb := db.Model(&DeviceInfo{})

	agencyId, err := modules.GetAgencyId2(ctx)
	if err != nil {
		return nil, err
	}
	if agencyId != 0 { // 是机构管理员的话，就需要添加机构排查
		tmpDb = tmpDb.Where("device_model in (?) AND agency_id=?", *batchUpdate.DeviceModels, agencyId)
	} else {
		tmpDb = tmpDb.Where("device_model in (?)", *batchUpdate.DeviceModels)
	}

	if batchUpdate.Tags != nil {
		tmpDb = tmpDb.Joins("JOIN tms_device_and_tag_mid b ON a.id = b.device_id AND b.deleted_at IS NULL AND b.tag_id IN (?)", *batchUpdate.Tags)
	}

	var ret []*DeviceInfo
	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return ret, err
	}

	return ret, nil
}
