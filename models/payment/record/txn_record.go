package record

import (
	"time"
	"tpayment/models"

	"github.com/shopspring/decimal"
)

type TxnRecord struct {
	models.BaseModel

	MerchantID            uint
	TotalAmount           decimal.Decimal
	Amount                decimal.Decimal
	Currency              string
	MerchantAccountID     uint
	TerminalID            string
	PaymentMethod         string
	PaymentEntryType      string
	PaymentType           string
	CustomerPaymentMethod string

	ReferenceID      string
	PartnerUUID      string
	AcquirerRRN      string
	AcquirerAuthCode string
	AcquirerReconID  string

	CompleteAt          time.Time
	AcquirerTxnDateTime time.Time

	Status               string
	VoidAt               time.Time
	RefundTimes          uint
	RefundAt             time.Time
	CaptureAt            time.Time
	GatewaySettlementAt  time.Time
	AcquirerSettlementAt time.Time

	PaymentFromType     string
	PaymentFromIP       string
	PaymentFromDeviceID string

	GatewayBatchNum  string
	AcquirerBatchNum string
	InvoiceNum       string

	ConsumerIdentify string // 消费者ID，信用卡用卡号，微信支付宝交易完成后会返回对应的ID
	CashierID        string // 操作人员ID
}
