package tms

import (
	"strconv"
	"strings"
	"tpayment/models"
	"tpayment/modules"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

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

func GetBatchUpdateRecordById(db *models.MyDB, ctx *gin.Context, id uint64) (*BatchUpdate, error) {
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

func QueryBatchUpdateRecord(db *models.MyDB, ctx *gin.Context, offset, limit uint64, filters map[string]string) (uint64, []*BatchUpdate, error) {
	agency := modules.IsAgencyAdmin(ctx)

	equalData := make(map[string]string)
	if agency != nil { // 是机构管理员的话，就需要添加机构排查
		equalData["agency_id"] = strconv.FormatUint(uint64(agency.ID), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Model(&BatchUpdate{}).Where(sqlCondition)

	// 统计总数
	var total uint64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*BatchUpdate
	if err = tmpDb.Order("updated_at desc").Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

//select * from tms_device where
//id in (select device_id from tms_device_and_tag_mid where tag_id IN (1,2) and deleted_at IS NULL group by device_id)
//and deleted_at is null;
func GetBatchUpdateDevices(db *models.MyDB, ctx *gin.Context, batchUpdate *BatchUpdate, offset uint64, limit uint64) ([]*DeviceInfo, error) {
	tmpDb := db.Model(&DeviceInfo{})

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
	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return ret, err
	}

	return ret, nil
}
