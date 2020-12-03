package payment_offline

import (
	"strconv"
	"time"
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/models/payment/acquirer"
	"tpayment/models/payment/merchantaccount"
	"tpayment/models/payment/paymentprocessrule"
	"tpayment/models/payment/record"
	"tpayment/pkg/paymentmethod/decodecardnum/creditcard"
	"tpayment/pkg/tlog"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// 预处理提交的数据，
// 1. 分析出支付方式，  2. 分析出用卡方式
func preHandleRequest(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
	var errCode conf.ResultCode

	// 匹配支付Payment Process Rule
	logger.Info("匹配支付Payment")
	errCode = decodePaymentData(ctx, txn)
	if errCode != conf.Success {
		return errCode
	}

	// 匹配payment process rule
	logger.Info("匹配payment process rule")
	if txn.OriginTxnID == 0 { // 首次交易
		logger.Info("首次交易")
		errCode = matchProcessRule2(ctx, txn)
		if errCode != conf.Success {
			return errCode
		}

		// 从payment process rule查找匹配merchant account和TID
		logger.Info("fetchMerchantAccountFirstTime")
		errCode = fetchMerchantAccountFirstTime(ctx, txn)
		if errCode != conf.Success {
			return errCode
		}
	} else {
		txn.PaymentProcessRule = new(paymentprocessrule.PaymentProcessRule)
		// 查找原始交易
		errCode = fetchOrgRecord(ctx, txn)
		if errCode != conf.Success {
			return errCode
		}

		// 从原始交易查找merchant account和TID
		errCode = fetchMerchantAccountFromOrg(ctx, txn)
		if errCode != conf.Success {
			return errCode
		}
	}

	// 提前生成record
	logger.Info("提前生成record")
	errCode = preBuildRecord(ctx, txn)
	if errCode != conf.Success {
		return errCode
	}

	return conf.Success
}

// 预处理Payment
func decodePaymentData(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	var err error
	// 分类支付方式
	switch txn.PaymentMethod {
	case conf.RequestCreditCard: // 常规信用卡
		txn.RealEntryType = txn.CreditCardBean.CardReaderMode
		if txn.AcquirerType == "" { // 没有收单行注明，则尝试解析
			txn.RealPaymentMethod, err = creditcard.Decode(txn.CreditCardBean.CardNumber)
			if err != nil {
				return conf.DecodeCardBrandError
			}
		} else {
			txn.RealPaymentMethod = txn.AcquirerType
		}
	}

	return conf.Success
}

// 匹配payment process rule
func matchProcessRule2(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)

	rule := new(paymentprocessrule.PaymentProcessRule)

	payRules, err := rule.GetByMerchantID(models.DB(), ctx, txn.MerchantID)
	if err != nil {
		logger.Error("GetByMerchantID fail->", err.Error())
		return conf.DBError
	}

	// 筛选出所有匹配支付方式，支付类型，用卡方式的payment rule
	var matchRules []*paymentprocessrule.PaymentProcessRule
	for i, payRule := range payRules {
		match := false
		// 匹配payment method
		if payRule.PaymentMethods != nil && len(*payRule.PaymentMethods) != 0 {
			for _, method := range *payRule.PaymentMethods {
				if method == txn.RealPaymentMethod {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		// 匹配entry type
		match = false
		if payRule.PaymentEntryTypes != nil && len(*payRule.PaymentEntryTypes) != 0 {
			for _, entryType := range *payRule.PaymentEntryTypes {
				if entryType == txn.RealEntryType {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		// 匹配支付方式
		match = false
		if payRule.PaymentTypes != nil && len(*payRule.PaymentTypes) != 0 {
			for _, paymentTypes := range *payRule.PaymentTypes {
				if paymentTypes == txn.TxnType {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		matchRules = append(matchRules, payRules[i])
	}

	// 未匹配到任何rule
	if len(matchRules) == 0 {
		return conf.NoPaymentProcessRule
	}

	txn.PaymentProcessRule = matchRules[0]

	return conf.Success
}

// 查找merchant account
func fetchMerchantAccountFirstTime(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)

	// 查找merchant account
	var err error
	merchantBean := &merchantaccount.MerchantAccount{
		BaseModel: models.BaseModel{
			Db:  models.DB(),
			Ctx: ctx,
		},
	}
	txn.PaymentProcessRule.MerchantAccount, err =
		merchantBean.Get(txn.PaymentProcessRule.MerchantAccountID)
	if err != nil {
		logger.Error("merchantBean.Get->", err.Error())
		return conf.DBError
	}
	if txn.PaymentProcessRule.MerchantAccount == nil {
		logger.Warn("can't find merchant account in payment process id->", txn.PaymentProcessRule.ID)
		return conf.ProcessRuleSettingError
	}

	// 分配TID
	if txn.AcquirerTerminalID != "" {
		tid := &acquirer.Terminal{
			BaseModel: models.BaseModel{
				Db:  models.DB(),
				Ctx: ctx,
			},
		}
		tid, err := tid.GetByTID(txn.PaymentProcessRule.MerchantAccountID, txn.AcquirerTerminalID)
		if err != nil {
			logger.Error("tid.GetByTID fail->", err.Error())
			return conf.DBError
		}

		if tid == nil {
			logger.Error("can't ->", err.Error())
			return conf.RecordNotFund
		}
		txn.PaymentProcessRule.MerchantAccount.Terminal = tid
	}

	return conf.Success
}

func fetchMerchantAccountFromOrg(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)

	txn.PaymentProcessRule = new(paymentprocessrule.PaymentProcessRule)
	// 查找merchant account
	var err error
	merchantBean := &merchantaccount.MerchantAccount{
		BaseModel: models.BaseModel{
			Db:  models.DB(),
			Ctx: ctx,
		},
	}
	txn.PaymentProcessRule.MerchantAccount, err =
		merchantBean.Get(txn.OrgRecord.MerchantAccountID)
	if err != nil {
		logger.Error("merchantBean.Get->", err.Error())
		return conf.DBError
	}
	if txn.PaymentProcessRule.MerchantAccount == nil {
		logger.Warn("can't find merchant account in payment process id->", txn.PaymentProcessRule.ID)
		return conf.ProcessRuleSettingError
	}

	// 查找acquirer
	acquirerBean := &agency.Acquirer{
		BaseModel: models.BaseModel{
			Db: models.DB(),
		},
	}
	txn.PaymentProcessRule.MerchantAccount.Acquirer, err =
		acquirerBean.Get(txn.PaymentProcessRule.MerchantAccount.AcquirerID)
	if err != nil {
		logger.Error("acquirerBean.Get->", err.Error())
		return conf.DBError
	}
	if txn.PaymentProcessRule.MerchantAccount.Acquirer == nil {
		logger.Warn("can't find acquirer in merchant id->", txn.PaymentProcessRule.MerchantAccount.ID)
		return conf.ProcessRuleSettingError
	}

	// 查找TID
	if txn.OrgRecord.TerminalID != 0 {
		terminalID := &acquirer.Terminal{
			BaseModel: models.BaseModel{
				Db:  models.DB(),
				Ctx: ctx,
			},
		}
		terminalID, err = terminalID.Get(txn.OrgRecord.TerminalID)
		if err != nil {
			logger.Error("terminalID.Get ", txn.OrgRecord.TerminalID, "->", err.Error())
			return conf.DBError
		}
		if terminalID == nil {
			logger.Error("can't find the terminal ", txn.OrgRecord.TerminalID)
			return conf.RecordNotFund
		}
		terminalID.BaseModel = models.BaseModel{
			Db:  models.DB(),
			Ctx: ctx,
		}
		txn.PaymentProcessRule.MerchantAccount.Terminal = terminalID
	}

	return conf.Success
}

// 查找原始交易
func fetchOrgRecord(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
	var err error

	recordBean := record.TxnRecord{
		BaseModel: models.BaseModel{
			Db:  models.DB(),
			Ctx: ctx,
		},
	}
	txn.OrgRecord, err = recordBean.GetByID(txn.OriginTxnID)
	if err != nil {
		logger.Warn("GetByID "+strconv.Itoa(int(txn.OriginTxnID))+" fail->", err.Error())
		return conf.DBError
	}
	if txn.OrgRecord == nil {
		logger.Warn("can't find the record " + strconv.Itoa(int(txn.OriginTxnID)))
		return conf.RecordNotFund
	}
	txn.OrgRecord.BaseModel = recordBean.BaseModel

	return conf.Success
}

// 生成交易记录
func preBuildRecord(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
	var (
		err    error
		amount decimal.Decimal
		ret    *record.TxnRecord
	)

	// 金额处理
	if txn.Amount != "" { // 没有传入金额
		amount, err = decimal.NewFromString(txn.Amount)
		if err != nil {
			logger.Warn("can't parse amount->", txn.Amount, ",", err.Error())
			return conf.ParameterError
		}
	} else {
		if txn.OrgRecord == nil {
			logger.Warn("can't get amount from org record->", txn.OriginTxnID)
			return conf.ParameterError
		}
		amount = txn.OrgRecord.Amount
	}

	if txn.OrgRecord == nil {
		ret = &record.TxnRecord{
			MerchantID:        txn.MerchantID,
			TotalAmount:       amount,
			Amount:            amount,
			Currency:          txn.Currency,
			MerchantAccountID: txn.PaymentProcessRule.MerchantID,
			PaymentMethod:     txn.RealPaymentMethod,
			PaymentEntryType:  txn.RealEntryType,
			ConsumerIdentify:  txn.CashierID,
		}
		if txn.CreditCardBean != nil && txn.CreditCardBean.CardNumber != "" {
			ret.ConsumerIdentify = txn.CreditCardBean.CardNumber
		}
	} else {
		ret = &record.TxnRecord{
			MerchantID:        txn.OrgRecord.MerchantID,
			TotalAmount:       txn.OrgRecord.Amount,
			Amount:            txn.OrgRecord.Amount,
			Currency:          txn.OrgRecord.Currency,
			MerchantAccountID: txn.OrgRecord.MerchantAccountID,
			PaymentMethod:     txn.OrgRecord.PaymentMethod,
			PaymentEntryType:  conf.ManualInput,
			ConsumerIdentify:  txn.OrgRecord.ConsumerIdentify,
		}
	}

	nowTime := time.Now()
	ret.IsOffline = true
	ret.PaymentFromName = txn.FromName
	ret.PaymentFromIP = txn.FromIp
	ret.PaymentFromDeviceID = txn.DeviceID
	ret.CompleteAt = &nowTime
	ret.AcquirerReconID = txn.AcquirerReconID
	ret.Status = record.Success
	ret.PartnerUUID = txn.Uuid
	ret.PaymentType = txn.TxnType

	txn.TxnRecord = ret
	return conf.Success
}
