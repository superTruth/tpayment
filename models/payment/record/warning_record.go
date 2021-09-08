package record

import (
	"tpayment/models"
)

type WarningRecord struct {
	models.BaseModel
	Uuid    string `json:"uuid"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (WarningRecord) TableName() string {
	return "payment_warning_record"
}

func (t *WarningRecord) Create(record *WarningRecord) error {
	return models.DB.Model(t).Create(record).Error
}

//// 生成异常交易
//func CreateWarningRecordByRequest(req *api_define.TxnReq, reason string) {
//	reqBytes, _ := json.Marshal(req)
//
//	wr := &WarningRecord{
//		BaseModel: models.BaseModel{
//			Db: models.DB(),
//		},
//		Uuid:    req.Uuid,
//		Reason:  reason,
//		Message: string(reqBytes),
//		Status:  "",
//	}
//
//	_ = wr.Create(wr)
//}
