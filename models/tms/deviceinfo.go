package tms

import (
	"strings"
	"tpayment/models"

	"github.com/jinzhu/gorm"
)

const (
	AppInDeviceExternalIdTypeDevice      = "merchantdevice"
	AppInDeviceExternalIdTypeBatchUpdate = "batch"
)

// devicetag
type DeviceTag struct {
	gorm.Model

	AgencyId string  `gorm:"column:agency_id"`
	Name     *string `gorm:"column:name"` // 外键
}

func (DeviceTag) TableName() string {
	return "mdm2_tags"
}

type DeviceTagFull struct {
	DeviceTag
	MidId uint `gorm:"column:mid_id"`
}

type DeviceAndTagMid struct {
	gorm.Model

	TagID    uint `gorm:"column:tag_id"`
	DeviceId uint `gorm:"column:device_id"`
}

func (DeviceAndTagMid) TableName() string {
	return "mdm2_device_and_tag_mid"
}

func GetDeviceByTag(tagsUuid []string, offset int, limit int) ([]string, error) {

	var deviceUuids []string

	// 全选
	if len(tagsUuid) == 0 {
		err := models.DB().Model(&DeviceAndTagMid{}).
			Offset(offset).Limit(offset).
			Select("device_sn").
			Find(deviceUuids).Error

		if err != nil {
			if gorm.ErrRecordNotFound == err { // 没有记录
				return deviceUuids, nil
			}
			return nil, err
		}

	}

	// 部分选择
	sb := strings.Builder{}
	for i, v := range tagsUuid {
		sb.WriteString(v)
		if i == len(tagsUuid)-1 { // 最后一个不用添加，
			break
		}
		sb.WriteString(",")
	}

	// 查找出所有device id
	rows, err := models.DB().Model(&DeviceAndTagMid{}).
		Where("tag_uuid in ( ? )", sb.String()).
		Select("device_sn").
		Offset(offset).Limit(limit).
		Rows()

	// nolint
	defer rows.Close()

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return deviceUuids, nil
		}
		return nil, err
	}

	err = rows.Scan(&deviceUuids)
	if err != nil {
		return nil, err
	}

	return deviceUuids, nil
}

// merchantdevice model
type DeviceModel struct {
	gorm.Model

	Name *string `gorm:"column:name"` // 外键
}

func (DeviceModel) TableName() string {
	return "mdm2_models"
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
