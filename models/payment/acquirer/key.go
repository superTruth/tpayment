package acquirer

import "tpayment/models"

type Key struct {
	models.BaseModel
	Tag   string `gorm:"column:tag"`
	Type  string `gorm:"column:type"`
	Value string `gorm:"column:value"`
}

func (Key) TableName() string {
	return "payment_acquirer_key"
}
