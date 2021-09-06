package merchant

import (
	"strconv"
	"tpayment/models"
	"tpayment/models/agency"

	"gorm.io/gorm"
)

var PaymentSettingDao = &PaymentSettingInDevice{}

type PaymentSettingInDevice struct {
	models.BaseModel

	MerchantDeviceId uint64              `json:"merchant_device_id,omitempty" gorm:"column:merchant_device_id"`
	PaymentMethods   *models.StringArray `json:"payment_methods,omitempty" gorm:"column:payment_methods;type:JSON"`
	EntryTypes       *models.StringArray `json:"entry_types,omitempty" gorm:"column:entry_types;type:JSON"`
	PaymentTypes     *models.StringArray `json:"payment_types,omitempty" gorm:"column:payment_types"`
	AcquirerId       uint64              `json:"acquirer_id,omitempty" gorm:"column:acquirer_id"`
	Mid              string              `json:"mid,omitempty" gorm:"column:mid"`
	Tid              string              `json:"tid,omitempty" gorm:"column:tid"`
	Addition         string              `json:"addition,omitempty" gorm:"column:addition"`
	AcquirerConfig   *agency.Acquirer    `json:"acquirer_config,omitempty" gorm:"-"`
}

func (PaymentSettingInDevice) TableName() string {
	return "merchant_payment_setting_in_device"
}

func GetPaymentSettingInDeviceById(id uint64) (*PaymentSettingInDevice, error) {
	ret := new(PaymentSettingInDevice)

	err := models.DB.Model(&PaymentSettingInDevice{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryPaymentSettingInDeviceRecord(merchantDeviceId, offset, limit uint64, filters map[string]string) (uint64, []*PaymentSettingInDevice, error) {
	var ret []*PaymentSettingInDevice

	equalData := make(map[string]string)
	equalData["merchant_device_id"] = strconv.FormatUint(merchantDeviceId, 10)
	sqlCondition := models.CombQueryCondition(equalData, filters)

	tmpDB := models.DB.Model(&PaymentSettingInDevice{}).Where(sqlCondition)

	// 统计总数
	var total int64 = 0
	err := tmpDB.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	// 查询记录
	err = tmpDB.Order("id desc").Offset(int(offset)).Limit(int(limit)).Find(&ret).Error
	if err != nil {
		return 0, nil, err
	}

	return uint64(total), ret, nil
}

func (p *PaymentSettingInDevice) GetByMidTid(deviceID uint64, mid, tid string) (*PaymentSettingInDevice, error) {
	ret := &PaymentSettingInDevice{}
	err := models.DB.Model(ret).
		Where("merchant_device_id=? and mid=? and tid=?", deviceID, mid, tid).
		First(ret).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
