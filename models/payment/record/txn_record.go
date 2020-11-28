package record

import (
	"time"
	"tpayment/models"

	"github.com/jinzhu/gorm"

	"github.com/shopspring/decimal"
)

const (
	Init         = "init"          // 初始化
	Success      = "success"       // 成功
	Fail         = "fail"          // 失败
	NeedReversal = "need_reversal" // 需要冲正
	Reversal     = "reversal"      // 冲正完成
	Pending      = "pending"       // 等待完成
)

type TxnRecord struct {
	models.BaseModel

	MerchantID            uint            `gorm:"column:merchant_id"`
	TotalAmount           decimal.Decimal `gorm:"column:total_amount"`
	Amount                decimal.Decimal `gorm:"column:amount"`
	Currency              string          `gorm:"column:currency"`
	MerchantAccountID     uint            `gorm:"column:merchant_account_id"`
	TerminalID            uint            `gorm:"column:terminal_id"`
	PaymentMethod         string          `gorm:"column:payment_method"`
	PaymentEntryType      string          `gorm:"column:payment_entry_type"`
	PaymentType           string          `gorm:"column:payment_type"`
	CustomerPaymentMethod string          `gorm:"column:customer_payment_method"`

	ReferenceID      string `gorm:"column:reference_id"`
	PartnerUUID      string `gorm:"column:partner_uuid"`
	AcquirerRRN      string `gorm:"column:acquirer_rrn"`
	AcquirerAuthCode string `gorm:"column:acquirer_auth_code"`
	AcquirerReconID  string `gorm:"column:acquirer_recon_id"`

	CompleteAt          *time.Time `gorm:"column:complete_at"`
	AcquirerTxnDateTime *time.Time `gorm:"column:acquirer_txn_date_time"`

	Status    string `gorm:"column:status"`
	ErrorCode string `json:"error_code"`
	ErrorDes  string `json:"error_des"`

	VoidAt               *time.Time `gorm:"column:void_at"`
	RefundTimes          uint       `gorm:"column:refund_times"`
	RefundAt             *time.Time `gorm:"column:refund_at"`
	CaptureAt            *time.Time `gorm:"column:capture_at"`
	GatewaySettlementAt  *time.Time `gorm:"column:gateway_settlement_at"`
	AcquirerSettlementAt *time.Time `gorm:"column:acquirer_settlement_at"`

	PaymentFromName     string `gorm:"column:payment_from_name"`
	PaymentFromIP       string `gorm:"column:payment_from_ip"`
	PaymentFromDeviceID string `gorm:"column:payment_from_device_id"`

	GatewayBatchNum  uint `gorm:"column:gateway_batch_num"`
	AcquirerBatchNum uint `gorm:"column:acquirer_batch_num"`
	InvoiceNum       uint `gorm:"column:invoice_num"`

	ConsumerIdentify string `gorm:"column:consumer_identify"` // 消费者ID，信用卡用卡号，微信支付宝交易完成后会返回对应的ID
	CashierID        string `gorm:"column:cashier_id"`        // 操作人员ID
}

func (TxnRecord) TableName() string {
	return "payment_txn_record"
}

func (t *TxnRecord) Create(record *TxnRecord) error {
	return t.Db.Model(t).Create(record).Error
}

func (t *TxnRecord) UpdateStatus(status string) error {
	err := t.Db.Model(t).Update(map[string]interface{}{"status": status}).Error
	if err != nil {
		return err
	}
	t.Status = status
	return nil
}

// 更新sale交易结果
func (t *TxnRecord) UpdateTxnResult() error {
	err := t.Db.Model(t).Where("id=?", t.ID).Select(
		[]string{"acquirer_rrn", "acquirer_auth_code",
			"acquirer_recon_id", "complete_at", "acquirer_txn_date_time",
			"status", "acquirer_batch_num", "consumer_identify"}).Updates(t).Error

	if err != nil {
		return err
	}
	return nil
}

// 更新void状态
func (t *TxnRecord) UpdateVoidStatus() error {
	err := t.Db.Model(t).Where("id=?", t.ID).Select(
		[]string{"void_at"}).Updates(t).Error

	if err != nil {
		return err
	}
	return nil
}

// 更新refund状态
func (t *TxnRecord) UpdateRefundStatus() error {
	err := t.Db.Model(t).Where("id=?", t.ID).Select(
		[]string{"total_amount", "refund_times", "refund_at"}).Updates(t).Error

	if err != nil {
		return err
	}
	return nil
}

// 查询一条记录
func (t *TxnRecord) GetByID(id uint) (*TxnRecord, error) {
	ret := new(TxnRecord)
	err := t.Db.Model(t).Where("id=?", id).First(ret).Error
	if gorm.ErrRecordNotFound == err { // 没有记录, 就创建一条记录
		if err != nil {
			return nil, err
		}
		return t, err
	}
	return ret, err
}
