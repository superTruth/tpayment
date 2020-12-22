package record

import "tpayment/models"

type AcquirerLog struct {
	models.BaseModel
	TxnID           uint64 `gorm:"column:txn_id"`
	RequestContent  string `gorm:"column:request_content"`
	ResponseContent string `gorm:"column:response_content"`
}

func (AcquirerLog) TableName() string {
	return "payment_acquirer_log"
}
