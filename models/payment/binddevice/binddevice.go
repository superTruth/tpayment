package binddevice

import "tpayment/models"

type BindDevice struct {
	models.BaseModel
	DeviceID                string `json:"device_id"`
	PaymentProcessingRuleID uint   `json:"payment_processing_rule_id"`
	TID                     string `json:"tid"`
	Addition                string `json:"addition"`
}
