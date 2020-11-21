package binddevice

import (
	"tpayment/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type BindDevice struct {
	models.BaseModel
	DeviceID                string `json:"device_id"`
	PaymentProcessingRuleID uint   `json:"payment_processing_rule_id"`
	TID                     string `json:"tid"`
	Addition                string `json:"addition"`
}

func (BindDevice) TableName() string {
	return "bind_device"
}

func (rule *BindDevice) Get(db *models.MyDB, ctx *gin.Context,
	processRuleID uint, deviceID string) (*BindDevice, error) {
	var ret *BindDevice
	err := db.Model(rule).Where("payment_processing_rule_id=? AND device_id=?",
		processRuleID, deviceID).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}
