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

	MerchantID            uint            `json:"merchant_id"`
	TotalAmount           decimal.Decimal `json:"total_amount"`
	Amount                decimal.Decimal `json:"amount"`
	Currency              string          `json:"currency"`
	MerchantAccountID     uint            `json:"merchant_account_id"`
	TerminalID            uint            `json:"terminal_id"`
	PaymentMethod         string          `json:"payment_method"`
	PaymentEntryType      string          `json:"payment_entry_type"`
	PaymentType           string          `json:"payment_type"`
	CustomerPaymentMethod string          `json:"customer_payment_method"`

	ReferenceID      string `json:"reference_id"`
	PartnerUUID      string `json:"partner_uuid"`
	AcquirerRRN      string `json:"acquirer_rrn"`
	AcquirerAuthCode string `json:"acquirer_auth_code"`
	AcquirerReconID  string `json:"acquirer_recon_id"`

	CompleteAt          *time.Time `json:"complete_at"`
	AcquirerTxnDateTime *time.Time `json:"acquirer_txn_date_time"`

	Status               string     `json:"status"`
	VoidAt               *time.Time `json:"void_at"`
	RefundTimes          uint       `json:"refund_times"`
	RefundAt             *time.Time `json:"refund_at"`
	CaptureAt            *time.Time `json:"capture_at"`
	GatewaySettlementAt  *time.Time `json:"gateway_settlement_at"`
	AcquirerSettlementAt *time.Time `json:"acquirer_settlement_at"`

	PaymentFromName     string `json:"payment_from_name"`
	PaymentFromIP       string `json:"payment_from_ip"`
	PaymentFromDeviceID string `json:"payment_from_device_id"`

	GatewayBatchNum  string `json:"gateway_batch_num"`
	AcquirerBatchNum string `json:"acquirer_batch_num"`
	InvoiceNum       string `json:"invoice_num"`

	ConsumerIdentify string `json:"consumer_identify"` // 消费者ID，信用卡用卡号，微信支付宝交易完成后会返回对应的ID
	CashierID        string `json:"cashier_id"`        // 操作人员ID
}

func (TxnRecord) TableName() string {
	return "txn_recode"
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
