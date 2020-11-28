package payment

import (
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
	logger := tlog.GetLogger(ctx)
	var errCode conf.ResultCode
	// 查找商户信息
	logger.Info("查找商户信息")
	errCode = fetchMerchantAgencyInfo2(ctx, txn)
	if errCode != conf.Success {
		return errCode
	}

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

// 获取商户和机构信息
func fetchMerchantAgencyInfo2(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
	var err error

	merchantBean := &merchant.Merchant{
		BaseModel: models.BaseModel{
			Db:  models.DB(),
			Ctx: ctx,
		},
	}
	merchantBean, err = merchantBean.Get(txn.MerchantID)
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
	agencyBean := &agency.Agency{
		BaseModel: models.BaseModel{
			Db:  models.DB(),
			Ctx: ctx,
		},
	}

	agencyBean, err = agencyBean.Get(merchantBean.AgencyId)
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
				CardExpDay:              applePayBean.ApplicationExpirationDate[4:],
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
	var errorCode conf.ResultCode

	// 查找merchant account
	var err error
	merchantBean := new(merchantaccount.MerchantAccount)
	txn.PaymentProcessRule.MerchantAccount, err =
		merchantBean.Get(models.DB(), ctx, txn.PaymentProcessRule.MerchantAccountID)
	if err != nil {
		logger.Error("merchantBean.Get->", err.Error())
		return conf.DBError
	}
	if txn.PaymentProcessRule.MerchantAccount == nil {
		logger.Warn("can't find merchant account in payment process id->", txn.PaymentProcessRule.ID)
		return conf.ProcessRuleSettingError
	}

	// 查找acquirer
	acquirerBean := new(agency.Acquirer)
	txn.PaymentProcessRule.MerchantAccount.Acquirer, err =
		acquirerBean.Get(models.DB(), ctx, txn.PaymentProcessRule.MerchantAccount.AcquirerID)
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
		tid := &acquirer.Terminal{
			BaseModel: models.BaseModel{
				Db:  models.DB(),
				Ctx: ctx,
			},
		}

		count, err := tid.GetTotal(txn.PaymentProcessRule.MerchantAccountID)
		if err != nil {
			logger.Warn("get total error->", err.Error())
			return conf.DBError
		}
		if count > 0 { // 没有tid，则认为是不需要绑定TID
			tid, errorCode =
				tid.GetOneAvailable(txn.PaymentProcessRule.MerchantAccountID, txn.DeviceID)
			if errorCode != conf.Success {
				logger.Warn("can't get available tid->", txn.DeviceID)
				return errorCode
			}
			txn.PaymentProcessRule.MerchantAccount.Terminal = tid // 可能会没有
			if txn.PaymentProcessRule.MerchantAccount.Terminal != nil {
				txn.PaymentProcessRule.MerchantAccount.Terminal.BaseModel.Db = models.DB()
				txn.PaymentProcessRule.MerchantAccount.Terminal.BaseModel.Ctx = ctx
			}
		}
	}

	return conf.Success
}

func fetchMerchantAccountFromOrg(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)

	txn.PaymentProcessRule = new(paymentprocessrule.PaymentProcessRule)
	// 查找merchant account
	var err error
	merchantBean := new(merchantaccount.MerchantAccount)
	txn.PaymentProcessRule.MerchantAccount, err =
		merchantBean.Get(models.DB(), ctx, txn.OrgRecord.MerchantAccountID)
	if err != nil {
		logger.Error("merchantBean.Get->", err.Error())
		return conf.DBError
	}
	if txn.PaymentProcessRule.MerchantAccount == nil {
		logger.Warn("can't find merchant account in payment process id->", txn.PaymentProcessRule.ID)
		return conf.ProcessRuleSettingError
	}

	// 查找acquirer
	acquirerBean := new(agency.Acquirer)
	txn.PaymentProcessRule.MerchantAccount.Acquirer, err =
		acquirerBean.Get(models.DB(), ctx, txn.PaymentProcessRule.MerchantAccount.AcquirerID)
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
			MerchantID:          txn.MerchantID,
			TotalAmount:         amount,
			Amount:              amount,
			Currency:            txn.Currency,
			MerchantAccountID:   txn.PaymentProcessRule.MerchantID,
			PaymentMethod:       txn.RealPaymentMethod,
			PaymentEntryType:    txn.RealEntryType,
			PaymentType:         txn.TxnType,
			PartnerUUID:         txn.Uuid,
			Status:              record.Init,
			PaymentFromName:     txn.FromName,
			PaymentFromIP:       txn.FromIp,
			PaymentFromDeviceID: txn.DeviceID,
			CashierID:           txn.CashierID,
		}
		if txn.CreditCardBean != nil && txn.CreditCardBean.CardNumber != "" {
			ret.ConsumerIdentify = txn.CreditCardBean.CardNumber
		}
	} else {
		ret = &record.TxnRecord{
			MerchantID:          txn.OrgRecord.MerchantID,
			TotalAmount:         txn.OrgRecord.Amount,
			Amount:              txn.OrgRecord.Amount,
			Currency:            txn.OrgRecord.Currency,
			MerchantAccountID:   txn.OrgRecord.MerchantAccountID,
			PaymentMethod:       txn.OrgRecord.PaymentMethod,
			PaymentEntryType:    conf.ManualInput,
			PaymentType:         txn.TxnType,
			PartnerUUID:         txn.Uuid,
			Status:              record.Init,
			PaymentFromName:     txn.FromName,
			PaymentFromIP:       txn.FromIp,
			PaymentFromDeviceID: txn.DeviceID,
			CashierID:           txn.CashierID,
			ConsumerIdentify:    txn.OrgRecord.ConsumerIdentify,
		}
	}

	txn.TxnRecord = ret
	return conf.Success
}
