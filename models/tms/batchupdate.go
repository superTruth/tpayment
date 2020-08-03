package tms

import (
	"tpayment/models"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type BatchUpdate struct {
	gorm.Model
	StoreID int `gorm:"column:store_id"`

	Description string `gorm:"column:description"`
	Status      int    `gorm:"column:status"`

	UpdateFailMsg string `gorm:"column:update_fail_msg"`

	Tags         string `gorm:"column:tags"`
	DeviceModels string `gorm:"column:device_models"`

	Apps []AppInDevice `gorm:"-"`
}

func (BatchUpdate) TableName() string {
	return "mdm2_batch_update"
}

func GetBatchUpdateRecord(db *models.MyDB, ctx echo.Context, id uint) (*BatchUpdate, error) {
	batchUpdate := new(BatchUpdate)

	err := models.DB().Where(&BatchUpdate{
		Model: gorm.Model{
			ID: id,
		},
	}).First(&batchUpdate).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}

		return nil, err
	}

	batchUpdate.Apps, _ = GetAppsInDevice(id, AppInDeviceExternalIdTypeBatchUpdate, 0, 1000) // 最多查出200条记录

	return batchUpdate, nil
}

func GetBatchUpdateDevices(batchUpdate *BatchUpdate, offset int, limit int) ([]DeviceInfo, error) {

	return nil, nil
	//sb := strings.Builder{}
	//
	//sb.WriteString("SELECT * FROM tms_device a ")
	//
	//tags, comErr := utils.JsonStringArray2StringArray(batchUpdate.Tags)
	//if comErr == nil && len(tags) != 0 { // 有选择tag的情况
	//	sb.WriteString("JOIN mdm2_device_and_tag_mid b ON a.id = b.device_id AND b.deleted_at IS NULL AND b.tag_id IN (")
	//	sb.WriteString(strings.Join(tags, ","))
	//	sb.WriteString(") ")
	//}
	//
	//deviceModels, comErr := utils.JsonStringArray2StringArray(batchUpdate.DeviceModels)
	//if comErr == nil && len(deviceModels) != 0 { // 有选择tag的情况
	//	sb.WriteString("JOIN mdm2_models c ON a.device_model = c.id AND c.deleted_at IS NULL AND c.id IN (")
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
