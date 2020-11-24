package acquirer

import "tpayment/models"

type Key struct {
	models.BaseModel
	Tag   string
	Type  string
	Value string
}

func (Key) TableName() string {
	return "acquirer_key"
}
