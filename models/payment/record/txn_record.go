package record

import (
	"errors"
	"time"
	"tpayment/conf"
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

var TxnRecordDao = &TxnRecord{}

type TxnRecord struct {
	models.BaseModel

	MerchantID  uint64          `gorm:"column:merchant_id" json:"merchant_id"`
	TotalAmount decimal.Decimal `gorm:"column:total_amount;type:decimal" json:"total_amount"`
	Amount      decimal.Decimal `gorm:"column:amount;type:decimal" json:"amount"`
	Currency    string          `gorm:"column:currency" json:"currency"`
	DccAmount   decimal.Decimal `gorm:"column:dcc_amount;type:decimal" json:"dcc_amount"`
	DccCurrency string          `gorm:"column:dcc_currency" json:"dcc_currency"`
	Installment uint64          `gorm:"column:installment" json:"installment"`

	MerchantAccountID     uint64 `gorm:"column:merchant_account_id" json:"merchant_account_id"`
	TerminalID            uint64 `gorm:"column:terminal_id" json:"terminal_id"`
	PaymentMethod         string `gorm:"column:payment_method" json:"payment_method"`
	PaymentEntryType      string `gorm:"column:payment_entry_type" json:"payment_entry_type"`
	PaymentType           string `gorm:"column:payment_type" json:"payment_type"`
	CustomerPaymentMethod string `gorm:"column:customer_payment_method" json:"customer_payment_method"`

	OrgTxnID         uint64 `gorm:"column:org_txn_id" json:"org_txn_id"`
	PartnerUUID      string `gorm:"column:partner_uuid" json:"partner_uuid"`
	AcquirerRRN      string `gorm:"column:acquirer_rrn" json:"acquirer_rrn"`
	AcquirerAuthCode string `gorm:"column:acquirer_auth_code" json:"acquirer_auth_code"`
	AcquirerReconID  string `gorm:"column:acquirer_recon_id" json:"acquirer_recon_id"`

	CompleteAt          *time.Time `gorm:"column:complete_at" json:"complete_at"`
	AcquirerTxnDateTime *time.Time `gorm:"column:acquirer_txn_date_time" json:"acquirer_txn_date_time"`

	Status    string `gorm:"column:status" json:"status"`
	ErrorCode string `gorm:"column:error_code" json:"error_code"`
	ErrorDes  string `gorm:"column:error_des" json:"error_des"`

	VoidAt               *time.Time `gorm:"column:void_at" json:"void_at"`
	RefundTimes          uint64     `gorm:"column:refund_times" json:"refund_times"`
	RefundAt             *time.Time `gorm:"column:refund_at" json:"refund_at"`
	AdjustTimes          uint64     `gorm:"column:adjust_times" json:"adjust_times"`
	AdjustAt             *time.Time `gorm:"column:adjust_at" json:"adjust_at"`
	CaptureAt            *time.Time `gorm:"column:capture_at" json:"capture_at"`
	IsOffline            bool       `gorm:"column:is_offline" json:"is_offline"`
	GatewaySettlementAt  *time.Time `gorm:"column:gateway_settlement_at" json:"gateway_settlement_at"`
	AcquirerSettlementAt *time.Time `gorm:"column:acquirer_settlement_at" json:"acquirer_settlement_at"`

	PaymentFromName     string `gorm:"column:payment_from_name" json:"payment_from_name"`
	PaymentFromIP       string `gorm:"column:payment_from_ip" json:"payment_from_ip"`
	PaymentFromDeviceID string `gorm:"column:payment_from_device_id" json:"payment_from_device_id"`

	GatewayBatchNum  uint64 `gorm:"column:gateway_batch_num" json:"gateway_batch_num"`
	AcquirerBatchNum uint64 `gorm:"column:acquirer_batch_num" json:"acquirer_batch_num"`
	AcquirerTraceNum uint64 `gorm:"column:acquirer_trace_num" json:"acquirer_trace_num"`
	InvoiceNum       uint64 `gorm:"column:invoice_num" json:"invoice_num"`

	ConsumerIdentify string `gorm:"column:consumer_identify" json:"consumer_identify"` // 消费者ID，信用卡用卡号，微信支付宝交易完成后会返回对应的ID
	CashierID        string `gorm:"column:cashier_id" json:"cashier_id"`               // 操作人员ID
}

func (TxnRecord) TableName() string {
	return "payment_txn_record"
}

func (t *TxnRecord) Create(record *TxnRecord) error {
	return models.DB().Model(t).Create(record).Error
}

func (t *TxnRecord) UpdateAll() error {
	return models.DB().Model(t).Save(t).Error
}

func (t *TxnRecord) UpdateStatus() error {
	dbTmp := models.DB().Model(t).Update(map[string]interface{}{
		"status":     t.Status,
		"error_code": t.ErrorCode,
		"error_des":  t.ErrorDes,
	})

	err := dbTmp.Error
	if err != nil {
		return err
	}

	if dbTmp.RowsAffected == 0 {
		return errors.New("no record updated")
	}

	return nil
}

// 更新sale交易结果
func (t *TxnRecord) UpdateTxnResult(db *gorm.DB) error {
	tmpDB := models.DB().DB
	if db != nil {
		tmpDB = db
	}
	dbTmp := tmpDB.Model(t).Where("id=?", t.ID).Select(
		[]string{"acquirer_rrn", "acquirer_auth_code",
			"acquirer_recon_id", "complete_at", "acquirer_txn_date_time",
			"status", "acquirer_batch_num", "consumer_identify", "error_code", "error_des"}).Updates(t)
	err := dbTmp.Error

	if err != nil {
		return err
	}

	if dbTmp.RowsAffected == 0 {
		return errors.New("no record updated")
	}

	return nil
}

// 更新void状态
func (t *TxnRecord) UpdateVoidStatus() error {
	dbTmp := models.DB().Model(t).Where("id=?", t.ID).Select(
		[]string{"void_at"}).Updates(t)

	err := dbTmp.Error
	if err != nil {
		return err
	}

	if dbTmp.RowsAffected == 0 {
		return errors.New("no record updated")
	}

	return nil
}

// 更新refund状态
func (t *TxnRecord) UpdateRefundStatus() error {
	dbTmp := models.DB().Model(t).Where("id=?", t.ID).Select(
		[]string{"total_amount", "refund_times", "refund_at"}).Updates(t)
	err := dbTmp.Error

	if err != nil {
		return err
	}

	if dbTmp.RowsAffected == 0 {
		return errors.New("no record updated")
	}
	return nil
}

// 查询一条记录
func (t *TxnRecord) GetByID(id uint64) (*TxnRecord, error) {
	ret := new(TxnRecord)
	err := models.DB().Model(t).Where("id=?", id).First(ret).Error
	if gorm.ErrRecordNotFound == err { // 没有记录, 就创建一条记录
		if err != nil {
			return nil, err
		}
		return t, err
	}
	return ret, err
}

// 通过ID或者Partner uuid查询一条交易记录
func (t *TxnRecord) GetByIDOrUuid(merchantID uint64, id uint64, partnerUuid string) (*TxnRecord, error) {
	ret := new(TxnRecord)
	var err error
	if id != 0 {
		err = models.DB().Model(t).Where("id=? AND merchant_id=?", id, merchantID).First(ret).Error
	} else {
		err = models.DB().Model(t).Where("partner_uuid=? AND merchant_id=?", partnerUuid, merchantID).First(ret).Error
	}

	if gorm.ErrRecordNotFound == err { // 没有记录, 就创建一条记录
		if err != nil {
			return nil, err
		}
		return t, err
	}
	return ret, err
}

type SettlementTotal struct {
	Currency     string          `gorm:"column:currency" json:"currency"`
	SaleAmount   decimal.Decimal `gorm:"column:sale_amount;type:decimal" json:"sale_amount"`
	SaleCount    uint64          `gorm:"column:sale_count" json:"sale_count"`
	RefundAmount decimal.Decimal `gorm:"column:refund_amount;type:decimal" json:"refund_amount"`
	RefundCount  uint64          `gorm:"column:refund_count" json:"refund_count"`
}

// 获取统计信息
func (t *TxnRecord) GetSettlementTotal(mid, tid, batchNum uint64) ([]*SettlementTotal, error) {
	var (
		err error
		ret []*SettlementTotal
	)

	// 先查看有多少种货币代码
	var currencies []*TxnRecord
	err = models.DB().Model(t).Where(
		"merchant_account_id=? and terminal_id=? and acquirer_batch_num=? and payment_type=? "+
			"and void_at is null and status=?",
		mid, tid, batchNum, conf.Sale, Success).Group("currency").Select("currency").Find(&currencies).Error
	if err != nil {
		return nil, err
	}
	if len(currencies) == 0 { // 没有交易
		return nil, nil
	}

	for i := range currencies {
		totalTmp := new(SettlementTotal)
		err = models.DB().Table(t.TableName()).Select("sum(amount) as sale_amount, count(*) as sale_count").Where(
			"merchant_account_id=? and terminal_id=? and acquirer_batch_num=? and payment_type=? "+
				"and void_at is null and status=? and currency=? and deleted_at is null",
			mid, tid, batchNum, conf.Sale, Success, currencies[i].Currency).Find(totalTmp).Error

		if err != nil {
			return nil, err
		}

		totalTmp.Currency = currencies[i].Currency
		ret = append(ret, totalTmp)
	}

	return ret, err
}

// 获取批上送记录
func (t *TxnRecord) GetBatchUploadRecords(mid, tid, batchNum, offset, limit uint64) ([]*TxnRecord, error) {
	var ret []*TxnRecord
	err := models.DB().Model(t).Where(
		"merchant_account_id=? and terminal_id=? and acquirer_batch_num=? and payment_type in (?) and void_at=0 and status=?",
		mid, tid, batchNum, []string{conf.Sale, conf.Refund, conf.PreAuthComplete}, Success).
		Offset(offset).Limit(limit).Find(&ret).Error

	if err != nil {
		return nil, err
	}

	return ret, err
}
