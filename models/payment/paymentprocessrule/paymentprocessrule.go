package paymentprocessrule

import "tpayment/models"

type PaymentProcessRule struct {
	models.BaseModel

	MerchantAccountID uint                `json:"merchant_account_id"`
	PaymentMethods    *models.StringArray `json:"payment_methods"`
	PaymentEntryTypes *models.StringArray `json:"payment_entry_types"`
	PaymentTypes      *models.StringArray `json:"payment_types"`
}
