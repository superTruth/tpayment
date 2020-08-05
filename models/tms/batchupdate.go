package tms

import (
	"strconv"
	"tpayment/models"
	"tpayment/modules"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type BatchUpdate struct {
	models.BaseModel

	AgencyId uint `gorm:"column:agency_id" json:"agency_id"`

	Description string `gorm:"column:description" json:"description"`
	Status      string `gorm:"column:status" json:"status"`

	UpdateFailMsg string `gorm:"column:update_fail_msg" json:"update_fail_msg"`

	Tags         string `gorm:"column:tags" json:"tags"`
	DeviceModels string `gorm:"column:device_models" json:"device_models"`

	Apps []*AppInDevice `gorm:"-"`
}

func (BatchUpdate) TableName() string {
	return "tms_batch_update"
}

func GetBatchUpdateRecordById(db *models.MyDB, ctx echo.Context, id uint) (*BatchUpdate, error) {
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

func QueryBatchUpdateRecord(db *models.MyDB, ctx echo.Context, offset, limit uint, filters map[string]string) (uint, []*BatchUpdate, error) {
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

func GetBatchUpdateDevices(db *models.MyDB, ctx echo.Context, batchUpdate *BatchUpdate, offset int, limit int) ([]*DeviceInfo, error) {

	return nil, nil
	//sb := strings.Builder{}
	//
	//sb.WriteString("SELECT * FROM tms_device a ")
	//
	//tags, comErr := utils.JsonStringArray2StringArray(batchUpdate.Tags)
	//if comErr == nil && len(tags) != 0 { // 有选择tag的情况
	//	sb.WriteString("JOIN tms_device_and_tag_mid b ON a.id = b.device_id AND b.deleted_at IS NULL AND b.tag_id IN (")
	//	sb.WriteString(strings.Join(tags, ","))
	//	sb.WriteString(") ")
	//}
	//
	//deviceModels, comErr := utils.JsonStringArray2StringArray(batchUpdate.DeviceModels)
	//if comErr == nil && len(deviceModels) != 0 { // 有选择tag的情况
	//	sb.WriteString("JOIN tms_model c ON a.device_model = c.id AND c.deleted_at IS NULL AND c.id IN (")
	//	sb.WriteString(strings.Join(deviceModels, ","))
	//	sb.WriteString(") ")
	//}
	//
	//sb.WriteString(fmt.Sprintf("WHERE a.store_id = %d and a.deleted_at IS NULL ", *batchUpdate.StoreID))
	//
	//sb.WriteString("GROUP BY a.id ")
	//
	//sb.WriteString("limit ")
	//sb.WriteString(strconv.Itoa(limit))
	//
	//sb.WriteString(" ")
	//
	//sb.WriteString(" offset ")
	//sb.WriteString(strconv.Itoa(offset))
	//
	////
	//rows, err := models.DB().Raw(sb.String()).Rows()
	//if err != nil {
	//	return nil, err
	//}
	//defer rows.Close()
	//
	//var ret []DeviceInfo
	//for rows.Next() {
	//	tmpDevice := new(DeviceInfo)
	//
	//	models.DB().ScanRows(rows, tmpDevice)
	//
	//	ret = append(ret, *tmpDevice)
	//}
	//
	//return ret, nil
}

func UpdateBatchUpdate(record *BatchUpdate) error {
	err := models.DB().Model(record).Updates(record).Error

	if err != nil {
		return err
	}

	return nil
}
