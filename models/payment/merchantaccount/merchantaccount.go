package merchantaccount

import "tpayment/models"

type MerchantAccount struct {
	models.BaseModel
	Name       string              `json:"name"`
	Currencies *models.StringArray `json:"currencies"`
	AcquirerID uint                `json:"acquirer_id"`
	MID        string              `json:"mid"`
	Addition   string              `json:"addition"`
	Disable    bool                `json:"disable"`
}
