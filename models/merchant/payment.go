package merchant

import (
	"strconv"
	"tpayment/models"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type PaymentSettingInDevice struct {
	models.BaseModel

	MerchantDeviceId uint                `json:"merchant_device_id,omitempty" gorm:"column:merchant_device_id"`
	PaymentMethods   *models.StringArray `json:"payment_methods,omitempty" gorm:"column:payment_methods;type:JSON"`
	EntryTypes       *models.StringArray `json:"entry_types,omitempty" gorm:"column:entry_types;type:JSON"`
	PaymentTypes     *models.StringArray `json:"payment_types,omitempty" gorm:"column:payment_types"`
	AcquirerId       uint                `json:"acquirer_id,omitempty" gorm:"column:acquirer_id"`
	Mid              string              `json:"mid,omitempty" gorm:"column:mid"`
	Tid              string              `json:"tid,omitempty" gorm:"column:tid"`
	Addition         string              `json:"addition,omitempty" gorm:"column:addition"`
}

func (PaymentSettingInDevice) TableName() string {
	return "payment_setting_in_device"
}

func GetPaymentSettingInDeviceById(db *models.MyDB, ctx echo.Context, id uint) (*PaymentSettingInDevice, error) {
	ret := new(PaymentSettingInDevice)

	err := db.Model(&PaymentSettingInDevice{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryPaymentSettingInDeviceRecord(db *models.MyDB, ctx echo.Context, merchantDeviceId, offset, limit uint, filters map[string]string) (uint, []*PaymentSettingInDevice, error) {
	var ret []*PaymentSettingInDevice

	equalData := make(map[string]string)
	equalData["merchant_device_id"] = strconv.FormatUint(uint64(merchantDeviceId), 10)
	sqlCondition := models.CombQueryCondition(equalData, filters)

	tmpDB := db.Model(&PaymentSettingInDevice{}).Where(sqlCondition)

	// 统计总数
	var total uint = 0
	err := tmpDB.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	// 查询记录
	err = tmpDB.Offset(offset).Limit(limit).Find(&ret).Error
	if err != nil {
		return 0, nil, err
	}

	return total, ret, nil
}
