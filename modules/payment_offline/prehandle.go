package payment_offline

import (
	"fmt"
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
	"tpayment/pkg/id"
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
	if txn.OrgTxnID == 0 { // 首次交易
		logger.Info("首次交易")
		if txn.PaymentMethod != conf.RequestOther { // 未知交易不需要匹配
			errCode = matchProcessRule(ctx, txn)
			if errCode != conf.Success {
				return errCode
			}

			// 从payment process rule查找匹配merchant account和TID
			logger.Info("fetchMerchantAccountFirstTime")
			errCode = fetchMerchantAccountFirstTime(ctx, txn)
			if errCode != conf.Success {
				return errCode
			}
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
	logger := tlog.GetLogger(ctx)
	var err error
	// 分类支付方式
	switch txn.PaymentMethod {
	case conf.RequestCreditCard: // 常规信用卡
		if txn.CreditCardBean == nil {
			logger.Warn("credit card txn but without credit card data")
			return conf.ParameterError
		}
		txn.RealEntryType = txn.CreditCardBean.CardReaderMode
		if txn.AcquirerType == "" { // 没有收单行注明，则尝试解析
			txn.RealPaymentMethod, err = creditcard.Decode(txn.CreditCardBean.CardNumber)
			if err != nil {
				return conf.DecodeCardBrandError
			}
		} else {
			txn.RealPaymentMethod = txn.AcquirerType
		}
	case conf.RequestOther:
		txn.RealPaymentMethod = conf.Other
	}

	return conf.Success
}

// 匹配payment process rule
func matchProcessRule(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
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
	if txn.DeviceID != "" {
		tid := &acquirer.Terminal{
			BaseModel: models.BaseModel{
				Db:  models.DB(),
				Ctx: ctx,
			},
		}
		tid, err := tid.GetByTID(txn.PaymentProcessRule.MerchantAccountID, txn.DeviceID)
		if err != nil {
			logger.Error("tid.GetByTID fail->", err.Error())
			return conf.DBError
		}

		if tid == nil {
			logger.Error("can't get any tid")
			return conf.RecordNotFund
		}
		txn.PaymentProcessRule.MerchantAccount.Terminal = tid
	}

	return conf.Success
}

func fetchMerchantAccountFromOrg(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)

	if txn.OrgRecord.PaymentMethod == conf.Other {
		return conf.Success
	}

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
	txn.OrgRecord, err = recordBean.GetByID(txn.OrgTxnID)
	if err != nil {
		logger.Warn("GetByID "+strconv.Itoa(int(txn.OrgTxnID))+" fail->", err.Error())
		return conf.DBError
	}
	if txn.OrgRecord == nil {
		logger.Warn("can't find the record " + strconv.Itoa(int(txn.OrgTxnID)))
		return conf.RecordNotFund
	}
	txn.OrgRecord.BaseModel.Db = recordBean.BaseModel.Db

	return conf.Success
}

// 生成交易记录
func preBuildRecord(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
	var (
		err      error
		amount   decimal.Decimal
		currency string
	)

	// 金额处理
	if txn.Amount != "" { // 没有传入金额
		amount, err = decimal.NewFromString(txn.Amount)
		if err != nil {
			logger.Warn("can't parse amount->", txn.Amount, ",", err.Error())
			return conf.ParameterError
		}
		currency = txn.Currency
	} else {
		if txn.OrgRecord == nil {
			logger.Warn("can't get amount from org record->", txn.OrgTxnID)
			return conf.ParameterError
		}
		amount = txn.OrgRecord.Amount
		currency = txn.OrgRecord.Currency
	}

	txn.TxnRecord = new(record.TxnRecord)
	txn.TxnRecordDetail = new(record.TxnRecordDetail)

	// 直接从request提取数据
	nowTime := time.Now()
	txn.TxnRecord.ID = id.New()
	txn.TxnRecord.Amount = amount
	txn.TxnRecord.Currency = currency
	txn.TxnRecord.TotalAmount = amount
	txn.TxnRecord.PartnerUUID = txn.Uuid
	txn.TxnRecord.PaymentType = txn.TxnType
	txn.TxnRecord.CompleteAt = &nowTime
	txn.TxnRecord.Status = record.Success
	txn.TxnRecord.IsOffline = true
	txn.TxnRecord.PaymentFromName = txn.FromName
	txn.TxnRecord.PaymentFromIP = txn.FromIp
	txn.TxnRecord.PaymentFromDeviceID = txn.DeviceID
	txn.TxnRecord.AcquirerRRN = txn.AcquirerRRN
	txn.TxnRecord.AcquirerReconID = txn.AcquirerReconID
	txn.TxnRecord.AcquirerTxnDateTime = txn.DateTime
	txn.TxnRecord.InvoiceNum = txn.InvoiceNum
	txn.TxnRecord.CashierID = txn.CashierID

	txn.TxnRecordDetail.ID = txn.TxnRecord.ID
	txn.TxnRecordDetail.Addition = txn.AdditionData

	// 区分第一次交易和第二次交易
	if txn.OrgRecord != nil {
		logger.Info("第二次交易")
		txn.TxnRecord.MerchantID = txn.OrgRecord.MerchantID
		txn.TxnRecord.MerchantAccountID = txn.OrgRecord.MerchantAccountID
		txn.TxnRecord.TerminalID = txn.OrgRecord.TerminalID
		txn.TxnRecord.PaymentMethod = txn.OrgRecord.PaymentMethod
		txn.TxnRecord.PaymentEntryType = conf.ManualInput
		txn.TxnRecord.CustomerPaymentMethod = txn.OrgRecord.CustomerPaymentMethod
		txn.TxnRecord.OrgTxnID = txn.OrgRecord.ID
	} else {
		logger.Info("第一次交易")
		if txn.PaymentProcessRule != nil && txn.PaymentProcessRule.MerchantAccount != nil {
			txn.TxnRecord.MerchantAccountID = txn.PaymentProcessRule.MerchantAccountID
			if txn.PaymentProcessRule.MerchantAccount.Terminal != nil {
				logger.Info("匹配terminal account")
				txn.TxnRecord.TerminalID = txn.PaymentProcessRule.MerchantAccount.Terminal.ID
			}
		}
		txn.TxnRecord.MerchantID = txn.MerchantID
		txn.TxnRecord.PaymentMethod = txn.RealPaymentMethod
		txn.TxnRecord.PaymentEntryType = txn.RealEntryType
		txn.TxnRecord.CustomerPaymentMethod = txn.CustomerPaymentMethod
	}

	// 特殊类型交易
	if txn.CreditCardBean != nil {
		txn.TxnRecord.AcquirerAuthCode = txn.CreditCardBean.AuthCode
		txn.TxnRecord.AcquirerBatchNum = txn.CreditCardBean.BatchNum
		txn.TxnRecord.AcquirerTraceNum = txn.CreditCardBean.TraceNum
		txn.TxnRecord.ConsumerIdentify = txn.CreditCardBean.CardNumber

		txn.TxnRecordDetail.CreditCardExp = fmt.Sprintf("%02s%02s%02s",
			txn.CreditCardBean.CardExpYear, txn.CreditCardBean.CardExpMonth, txn.CreditCardBean.CardExpDay)
		txn.TxnRecordDetail.CreditCardFallBack = txn.CreditCardBean.CardFallback
		txn.TxnRecordDetail.CreditCardSN = txn.CreditCardBean.CardSn
		txn.TxnRecordDetail.CreditCardHolderName = txn.CreditCardBean.CardHolderName
		txn.TxnRecordDetail.CreditCardIsMsdCard = txn.CreditCardBean.IsMsdCard
		txn.TxnRecordDetail.CreditCardIccRequest = txn.CreditCardBean.IccRequest
		txn.TxnRecordDetail.CreditCardIccResponse = txn.CreditCardBean.IccResponse
		txn.TxnRecordDetail.ResponseCode = txn.CreditCardBean.ResponseCode
	}

	return conf.Success
}
