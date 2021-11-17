package pay_manage

type CheckRequest struct {
	MerchantId  uint64 `json:"merchant_id,omitempty"`
	TxnID       uint64 `json:"txn_id"`
	PartnerUUID string `json:"partner_uuid"`
}
