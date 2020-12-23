package pay_online

import (
	"fmt"
	"strconv"
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/models/merchant"
	"tpayment/models/payment/acquirer"
	"tpayment/models/payment/key"
	"tpayment/models/payment/merchantaccount"
	"tpayment/models/payment/paymentprocessrule"
	"tpayment/models/payment/record"
	"tpayment/pkg/id"
	"tpayment/pkg/paymentmethod/decodecardnum/applepay"
	"tpayment/pkg/paymentmethod/decodecardnum/creditcard"
	"tpayment/pkg/paymentmethod/decodecardnum/qrcode"
	"tpayment/pkg/tlog"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// 预处理提交的数据，
// 1. 分析出支付方式，  2. 分析出用卡方式
func preHandleRequest(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	var errCode conf.ResultCode
	// 查找商户信息
	errCode = fetchMerchantAgencyInfo2(ctx, txn)
	if errCode != conf.Success {
		return errCode
	}

	// 匹配支付Payment Process Rule
	errCode = decodePaymentData(ctx, txn)
	if errCode != conf.Success {
		return errCode
	}

	// 匹配payment process rule
	if txn.OrgTxnID == 0 { // 首次交易
		errCode = matchProcessRule2(ctx, txn)
		if errCode != conf.Success {
			return errCode
		}

		// 从payment process rule查找匹配merchant account和TID
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
	errCode = preBuildRecord2(ctx, txn)
	if errCode != conf.Success {
		return errCode
	}

	return conf.Success
}

// 获取商户和机构信息
func fetchMerchantAgencyInfo2(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
	var err error

	merchantBean, err := merchant.Dao.Get(txn.MerchantID)
	if err != nil {
		logger.Warn("merchant fetch fail->", err.Error())
		return conf.DBError
	}

	if merchantBean == nil {
		logger.Warn("can't find merchant")
		return conf.ProcessRuleSettingError
	}
	txn.MerchantInfo = merchantBean

	// Agency
	agencyBean, err := agency.Dao.Get(merchantBean.AgencyId)
	if err != nil {
		logger.Warn("agency fetch fail->", err.Error())
		return conf.DBError
	}

	if agencyBean == nil {
		logger.Warn("can't find agency")
		return conf.ProcessRuleSettingError
	}
	txn.AgencyInfo = agencyBean

	return conf.Success
}

// 预处理Payment
func decodePaymentData(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
	var err error

	// 分类支付方式
	switch txn.PaymentMethod {
	case conf.RequestCreditCard: // 常规信用卡
		txn.RealEntryType = txn.CreditCardBean.CardReaderMode
		txn.RealPaymentMethod, err = creditcard.Decode(txn.CreditCardBean.CardNumber)
		if err != nil {
			logger.Warn("creditcard.Decode error->", err.Error())
			return conf.DecodeCardBrandError
		}
	case conf.RequestApplePay: // apple pay
		txn.RealEntryType = conf.ApplePay

		// 查找配置的Key
		pukHash, err := applepay.GetApplePayKeyHash(txn.ApplePayBean.Token)
		if err != nil {
			logger.Warn("get apple pay hash fail->", err.Error())
			return conf.ParameterError
		}

		applePayKey := new(key.ApplePayKey)

		applePayKey, err = applePayKey.GetKeyByHash(models.DB(), ctx, pukHash)
		if err != nil {
			logger.Error("GetKeyByHash fail->", err.Error())
			return conf.DBError
		}
		if applePayKey == nil {
			logger.Info("use apple pay key->", applePayKey.ID)
			// 解码apple pay数据
			applePayBean, err := applepay.DecodeApplePay(txn.ApplePayBean.Token, &applepay.ConfigKey{
				PublicKey:  applePayKey.PublicKey,
				PrivateKey: applePayKey.PrivateKey,
			})

			if err != nil {
				logger.Warn("applepay.DecodeApplePay fail->", err.Error())
				return conf.DecodeError
			}
			txn.CreditCardBean = &api_define.CreditCardBean{
				CardExpMonth:            applePayBean.ApplicationExpirationDate[2:4],
				CardExpYear:             applePayBean.ApplicationExpirationDate[:2],
				CardNumber:              applePayBean.ApplicationPrimaryAccountNumber,
				CardHolderName:          applePayBean.CardholderName,
				IccRequest:              applePayBean.PaymentData.EmvData,
				PIN:                     applePayBean.PaymentData.EncryptedPINData,
				ECI:                     applePayBean.PaymentData.EciIndicator,
				OnlinePaymentCryptogram: applePayBean.PaymentData.OnlinePaymentCryptogram,
			}
		} else {
			logger.Info("txn can't find apple pay key")
		}
	case conf.RequestConsumerPresentQR: // 商户扫手机二维码
		txn.RealEntryType = conf.ConsumerPresentQR

		// 优先判断微信，支付宝，国内云闪付版本
		txn.RealPaymentMethod, err = qrcode.Decode(txn.ConsumerPresentQR.Content)
		if err == nil {
			break
		}

		// EMV 银行卡二维码
		emvQRContent, err := qrcode.DecodeEmvQR(txn.ConsumerPresentQR.Content)
		if err == nil {
			txn.RealPaymentMethod, err = creditcard.Decode(emvQRContent.CardNum)
			if err != nil {
				logger.Warn("creditcard.Decode error->", err.Error())
				return conf.DecodeCardBrandError
			}

			txn.CreditCardBean = &api_define.CreditCardBean{
				CardNumber: emvQRContent.CardNum,
				CardTrack2: emvQRContent.Track2,
				CardSn:     emvQRContent.CardSn,
				IccRequest: emvQRContent.ICCData,
			}

			txn.RealPaymentMethod, err = creditcard.Decode(txn.CreditCardBean.CardNumber)
			if err != nil {
				logger.Warn("creditcard.Decode error->", err.Error())
				return conf.DecodeCardBrandError
			}
			return conf.Success
		}

		logger.Warn("can't decode the qr code->", txn.ConsumerPresentQR.Content)
		return conf.DecodeQRError
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
	txn.PaymentProcessRule.MerchantAccount, err =
		merchantaccount.Dao.Get(txn.PaymentProcessRule.MerchantAccountID)
	if err != nil {
		logger.Error("merchantBean.Get->", err.Error())
		return conf.DBError
	}
	if txn.PaymentProcessRule.MerchantAccount == nil {
		logger.Warn("can't find merchant account in payment process id->", txn.PaymentProcessRule.ID)
		return conf.ProcessRuleSettingError
	}

	// 查找acquirer
	txn.PaymentProcessRule.MerchantAccount.Acquirer, err =
		agency.AcquirerDao.Get(txn.PaymentProcessRule.MerchantAccount.AcquirerID)
	if err != nil {
		logger.Error("acquirerBean.Get->", err.Error())
		return conf.DBError
	}
	if txn.PaymentProcessRule.MerchantAccount.Acquirer == nil {
		logger.Warn("can't find acquirer in merchant id->", txn.PaymentProcessRule.MerchantAccount.ID)
		return conf.ProcessRuleSettingError
	}

	// 分配TID
	if txn.DeviceID != "" {
		count, err := acquirer.TerminalDao.GetTotal(txn.PaymentProcessRule.MerchantAccountID)
		if err != nil {
			logger.Warn("get total error->", err.Error())
			return conf.DBError
		}
		if count > 0 { // 没有tid，则认为是不需要绑定TID
			tid, errorCode :=
				acquirer.TerminalDao.GetOneAvailable(txn.PaymentProcessRule.MerchantAccountID, txn.DeviceID)
			if errorCode != conf.Success {
				logger.Warn("can't get available tid->", txn.DeviceID)
				return errorCode
			}
			txn.PaymentProcessRule.MerchantAccount.Terminal = tid // 可能会没有
		}
	}

	return conf.Success
}

func fetchMerchantAccountFromOrg(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)

	txn.PaymentProcessRule = new(paymentprocessrule.PaymentProcessRule)
	// 查找merchant account
	var err error
	txn.PaymentProcessRule.MerchantAccount, err =
		merchantaccount.Dao.Get(txn.OrgRecord.MerchantAccountID)
	if err != nil {
		logger.Error("merchantBean.Get->", err.Error())
		return conf.DBError
	}
	if txn.PaymentProcessRule.MerchantAccount == nil {
		logger.Warn("can't find merchant account in payment process id->", txn.PaymentProcessRule.ID)
		return conf.ProcessRuleSettingError
	}

	// 查找acquirer
	txn.PaymentProcessRule.MerchantAccount.Acquirer, err =
		agency.AcquirerDao.Get(txn.PaymentProcessRule.MerchantAccount.AcquirerID)
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
		terminalID, err := acquirer.TerminalDao.Get(txn.OrgRecord.TerminalID)
		if err != nil {
			logger.Error("terminalID.Get ", txn.OrgRecord.TerminalID, "->", err.Error())
			return conf.DBError
		}
		if terminalID == nil {
			logger.Error("can't find the terminal ", txn.OrgRecord.TerminalID)
			return conf.RecordNotFund
		}
		txn.PaymentProcessRule.MerchantAccount.Terminal = terminalID
	}

	return conf.Success
}

// 查找原始交易
func fetchOrgRecord(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
	var err error

	txn.OrgRecord, err = record.TxnRecordDao.GetByID(txn.OrgTxnID)
	if err != nil {
		logger.Warn("GetByID "+strconv.Itoa(int(txn.OrgTxnID))+" fail->", err.Error())
		return conf.DBError
	}
	if txn.OrgRecord == nil {
		logger.Warn("can't find the record " + strconv.Itoa(int(txn.OrgTxnID)))
		return conf.RecordNotFund
	}

	return conf.Success
}

// 生成交易记录
func preBuildRecord2(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
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
	txn.TxnRecord.ID = id.New()
	txn.TxnRecord.Amount = amount
	txn.TxnRecord.Currency = currency
	txn.TxnRecord.TotalAmount = amount
	txn.TxnRecord.PartnerUUID = txn.Uuid
	txn.TxnRecord.PaymentType = txn.TxnType
	txn.TxnRecord.Status = record.Init
	txn.TxnRecord.IsOffline = false
	txn.TxnRecord.PaymentFromName = txn.FromName
	txn.TxnRecord.PaymentFromIP = txn.FromIp
	txn.TxnRecord.PaymentFromDeviceID = txn.DeviceID
	txn.TxnRecord.InvoiceNum = txn.InvoiceNum
	txn.TxnRecord.CashierID = txn.CashierID

	txn.TxnRecordDetail.ID = txn.TxnRecord.ID
	txn.TxnRecordDetail.TxnExpAt = txn.TxnExpAt
	txn.TxnRecordDetail.RedirectSuccessUrl = txn.RedirectSuccessUrl
	txn.TxnRecordDetail.RedirectFailUrl = txn.RedirectFailUrl
	txn.TxnRecordDetail.ResultNotifyUrl = txn.ResultNotifyUrl

	// 区分第一次交易和第二次交易
	if txn.OrgRecord != nil {
		txn.TxnRecord.MerchantID = txn.OrgRecord.MerchantID
		txn.TxnRecord.MerchantAccountID = txn.OrgRecord.MerchantAccountID
		txn.TxnRecord.TerminalID = txn.OrgRecord.TerminalID
		txn.TxnRecord.PaymentMethod = txn.OrgRecord.PaymentMethod
		txn.TxnRecord.PaymentEntryType = conf.ManualInput
		txn.TxnRecord.CustomerPaymentMethod = txn.OrgRecord.CustomerPaymentMethod
		txn.TxnRecord.OrgTxnID = txn.OrgRecord.ID

		if txn.OrgRecordDetail != nil {
			txn.TxnRecordDetail.CreditCardExp = txn.OrgRecordDetail.CreditCardExp
			txn.TxnRecordDetail.CreditCardHolderName = txn.OrgRecordDetail.CreditCardHolderName
		}

	} else {
		if txn.PaymentProcessRule != nil && txn.PaymentProcessRule.MerchantAccount != nil {
			txn.TxnRecord.MerchantAccountID = txn.PaymentProcessRule.MerchantAccountID
			if txn.PaymentProcessRule.MerchantAccount.Terminal != nil {
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
		txn.TxnRecord.AcquirerBatchNum = txn.CreditCardBean.BatchNum
		txn.TxnRecord.AcquirerTraceNum = txn.CreditCardBean.TraceNum
		txn.TxnRecord.ConsumerIdentify = txn.CreditCardBean.CardNumber

		txn.TxnRecordDetail.CreditCardExp = fmt.Sprintf("%02s%02s",
			txn.CreditCardBean.CardExpYear, txn.CreditCardBean.CardExpMonth)
		txn.TxnRecordDetail.CreditCardFallBack = txn.CreditCardBean.CardFallback
		txn.TxnRecordDetail.CreditCardSN = txn.CreditCardBean.CardSn
		txn.TxnRecordDetail.CreditCardHolderName = txn.CreditCardBean.CardHolderName
		txn.TxnRecordDetail.CreditCardIsMsdCard = txn.CreditCardBean.IsMsdCard
		txn.TxnRecordDetail.CreditCardIccRequest = txn.CreditCardBean.IccRequest
		txn.TxnRecordDetail.CreditCardIccResponse = txn.CreditCardBean.IccResponse
		txn.TxnRecordDetail.ResponseCode = txn.CreditCardBean.ResponseCode
	}

	// ApplePayBean交易
	if txn.CreditCardTokenBean != nil {
		txn.TxnRecordDetail.Token = txn.CreditCardTokenBean.Token
	}

	// 3DS交易
	if txn.CreditCard3DSBean != nil {
		txn.TxnRecordDetail.TDSEnable = true
	}

	// 扫码交易
	if txn.ConsumerPresentQR != nil {
		txn.TxnRecordDetail.Token = txn.ConsumerPresentQR.Content
		txn.TxnRecordDetail.TokenType = txn.ConsumerPresentQR.CodeType
	}

	return conf.Success
}
