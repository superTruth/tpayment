package tms

import (
	"strconv"
	"strings"
	"tpayment/models"
	"tpayment/modules"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var BatchUpdateDao = &BatchUpdate{}

type BatchUpdate struct {
	models.BaseModel

	AgencyId uint64 `gorm:"column:agency_id" json:"agency_id"`

	Description string `gorm:"column:description" json:"description"`
	Status      string `gorm:"column:status" json:"status"`

	UpdateFailMsg string `gorm:"column:update_fail_msg" json:"update_fail_msg"`

	Tags         *models.IntArray `gorm:"column:tags" json:"-"`
	DeviceModels *models.IntArray `gorm:"column:device_models" json:"-"`

	ConfigTags []*DeviceTag `gorm:"-" json:"tags"`

	ConfigModels []*DeviceModel `gorm:"-" json:"device_models"`

	Apps []*AppInDevice `gorm:"-" json:"-"`
}

func (BatchUpdate) TableName() string {
	return "tms_batch_update"
}

func GetBatchUpdateRecordById(id uint64) (*BatchUpdate, error) {
	ret := new(BatchUpdate)

	err := models.DB.Model(&BatchUpdate{}).Where("id=?", id).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}

		return nil, err
	}

	return ret, nil
}

func QueryBatchUpdateRecord(ctx *gin.Context, offset, limit uint64, filters map[string]string) (uint64, []*BatchUpdate, error) {
	agency := modules.IsAgencyAdmin(ctx)

	equalData := make(map[string]string)
	if agency != nil { // 是机构管理员的话，就需要添加机构排查
		equalData["agency_id"] = strconv.FormatUint(uint64(agency.ID), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := models.DB.Model(&BatchUpdate{}).Where(sqlCondition)

	// 统计总数
	var total int64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*BatchUpdate
	if err = tmpDb.Order("id desc").Offset(int(offset)).Limit(int(limit)).Find(&ret).Error; err != nil {
		return uint64(total), ret, err
	}

	return uint64(total), ret, nil
}

//select * from tms_device where
//id in (select device_id from tms_device_and_tag_mid where tag_id IN (1,2) and deleted_at IS NULL group by device_id)
//and deleted_at is null;
func GetBatchUpdateDevices(ctx *gin.Context, batchUpdate *BatchUpdate, offset uint64, limit uint64) ([]*DeviceInfo, error) {
	tmpDb := models.DB.Model(&DeviceInfo{})

	agencyId, err := modules.GetAgencyId2(ctx)
	if err != nil {
		return nil, err
	}

	midSqlSb := strings.Builder{}
	// mid sql 中间表创建
	// tags id
	tagsSql := ""
	if batchUpdate.Tags != nil && len(*batchUpdate.Tags) != 0 {
		tagsSql = "tag_id IN (" + batchUpdate.Tags.String() + ") "
	}
	if tagsSql != "" {
		midSqlSb.WriteString("select device_id from tms_device_and_tag_mid where ")
		midSqlSb.WriteString(tagsSql)
		midSqlSb.WriteString("and deleted_at IS NULL group by device_id")
	}

	// 组建整体的Sql
	retSqlSb := strings.Builder{}
	if midSqlSb.Len() != 0 {
		retSqlSb.WriteString("id in (")
		retSqlSb.WriteString(midSqlSb.String())
		retSqlSb.WriteString(") ")
	}

	// model id
	if batchUpdate.DeviceModels != nil && len(*batchUpdate.DeviceModels) != 0 {
		if retSqlSb.Len() != 0 {
			retSqlSb.WriteString(" AND ")
		}
		retSqlSb.WriteString("device_model in (" + batchUpdate.DeviceModels.String() + ") ")
	}

	// agency filter
	if agencyId != 0 {
		if retSqlSb.Len() != 0 {
			retSqlSb.WriteString(" AND ")
		}
		retSqlSb.WriteString("agency_id=")
		retSqlSb.WriteString(strconv.FormatUint(uint64(agencyId), 10))
	}

	tmpDb = tmpDb.Where(retSqlSb.String())

	//
	//tmpSql := ""
	//if agencyId != 0 { // 是机构管理员的话，就需要添加机构排查
	//  if batchUpdate.DeviceModels != nil && len(*batchUpdate.DeviceModels) != 0 {
	//      tmpSql = "device_model in (" + batchUpdate.DeviceModels.String() + ") AND "
	//  }
	//  tmpDb = tmpDb.Where(tmpSql+" agency_id=?", agencyId)
	//} else {
	//  if batchUpdate.DeviceModels != nil && len(*batchUpdate.DeviceModels) != 0 {
	//      tmpSql = "device_model in (" + batchUpdate.DeviceModels.String() + ")"
	//  }
	//  tmpDb = tmpDb.Where(tmpSql)
	//}
	//
	//if batchUpdate.Tags != nil {
	//  if batchUpdate.Tags != nil && len(*batchUpdate.Tags) != 0 {
	//      tmpSql = "AND b.tag_id IN (" + batchUpdate.Tags.String() + ")"
	//  } else {
	//      tmpSql = ""
	//  }
	//
	//  tmpDb = tmpDb.Joins("JOIN tms_device_and_tag_mid b ON tms_device.id = b.device_id AND b.deleted_at IS NULL " + tmpSql)
	//}

	var ret []*DeviceInfo
	if err = tmpDb.Offset(int(offset)).Limit(int(limit)).Find(&ret).Error; err != nil {
		return ret, err
	}

	return ret, nil
}

const (
	BatchUpdateStatusPending = "pending"
	BatchUpdateStatusSuccess = "success"
)

var DeviceInBatchDao = &DeviceInBatchUpdate{}

type DeviceInBatchUpdate struct {
	models.BaseModel

	BatchID  uint64 `gorm:"column:batch_id" json:"-"`
	DeviceID uint64 `gorm:"column:device_id" json:"-"`
	Status   string `gorm:"column:status" json:"-"`

	DeviceInfo *DeviceInfo `gorm:"-" json:"device_info"`

	BatchUpdate *BatchUpdate `gorm:"-" json:"-"`
}

func (DeviceInBatchUpdate) TableName() string {
	return "tms_device_batch_update"
}

func (d *DeviceInBatchUpdate) GetDevicesByBatch(batchID, offset, limit uint64) (uint64, []*DeviceInBatchUpdate, error) {
	var ret []*DeviceInBatchUpdate

	tmpDb := models.DB.Model(&DeviceInBatchUpdate{}).
		Where("batch_id=?", batchID)
	// 统计总数
	var total int64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	err = tmpDb.Offset(int(offset)).Limit(int(limit)).Find(&ret).Error
	if err != nil {
		return 0, nil, err
	}

	for i := 0; i < len(ret); i++ {
		ret[i].DeviceInfo, err = DeviceInfoDao.GetByID(ret[i].DeviceID)
		if err != nil {
			return 0, nil, err
		}
	}

	return uint64(total), ret, nil
}

func (d *DeviceInBatchUpdate) Create(data *DeviceInBatchUpdate) error {
	return models.DB.Create(data).Error
}

func (d *DeviceInBatchUpdate) GetUnCompletedBatchByDevice(deviceID uint64) ([]*DeviceInBatchUpdate, error) {
	var ret []*DeviceInBatchUpdate
	err := models.DB.Model(&DeviceInBatchUpdate{}).
		Where("device_id=? and status != ?", deviceID, BatchUpdateStatusSuccess).Find(&ret).Error
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (d *DeviceInBatchUpdate) UpdateStatus(data *DeviceInBatchUpdate) error {
	return models.DB.Model(d).Select("status").Updates(data).Error
}

func (d *DeviceInBatchUpdate) GetByBatchIDDeviceID(batchID, deviceID uint64) (*DeviceInBatchUpdate, error) {
	ret := &DeviceInBatchUpdate{}
	err := models.DB.Model(ret).
		Where("batch_id=? AND device_id=?", batchID, deviceID).
		First(ret).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
