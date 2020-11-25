package payment

import (
	"errors"
	"time"
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/models/merchant"
	"tpayment/models/payment/record"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

const SaleMaxExpTime = time.Minute * 5

func SaleHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(api_define.TxnReq)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 获取机构/商户信息
	merchantInfo, agencyInfo, errorCode := fetchMerchantAgencyInfo(ctx, req.MerchantID)
	if errorCode != conf.Success {
		logger.Warn("fetchMerchantAgencyInfo fail->", errorCode.String())
		modules.BaseError(ctx, errorCode)
		return
	}

	// 预处理交易数据，分析出真实卡号，用卡方式
	errorCode = preHandleRequest(ctx, req)
	if errorCode != conf.Success {
		logger.Warn("preHandleRequest fail->", errorCode.String())
		modules.BaseError(ctx, errorCode)
		return
	}

	// 提取payment processing rule
	req.PaymentProcessRule, errorCode = matchProcessRule(ctx, req)
	if errorCode != conf.Success {
		logger.Warn("matchProcessRule fail->", errorCode.String())
		modules.BaseError(ctx, errorCode)
		return
	}
	if req.PaymentProcessRule.MerchantAccount.TID != nil { // 如果有TID的情况，需要锁定TID
		errorCode = req.PaymentProcessRule.MerchantAccount.TID.Lock(SaleMaxExpTime)
		if errorCode != conf.Success {
			modules.BaseError(ctx, errorCode)
			return
		}
		defer req.PaymentProcessRule.MerchantAccount.TID.UnLock()
	}

	// 获取merchant account, acquirer
	errorCode = fetchMerchantAccount(ctx, req)
	if errorCode != conf.Success {
		logger.Warn("fetchMerchantAccount fail->", errorCode.String())
		modules.BaseError(ctx, errorCode)
		return
	}

	// 获取sale交易对象
	acquirerImpl, ok := acquirer_impl.AcquirerImpls[req.PaymentProcessRule.MerchantAccount.Acquirer.Name]
	if !ok {
		logger.Warn("can't find acquirer impl->", req.PaymentProcessRule.MerchantAccount.Acquirer.Name)
		modules.BaseError(ctx, conf.UnknownError)
		return
	}
	saleImp, ok := acquirerImpl.(acquirer_impl.ISale)
	if !ok {
		logger.Warn("the acquirer not support sale->", req.PaymentProcessRule.MerchantAccount.Acquirer.Name)
		modules.BaseError(ctx, conf.UnknownError)
		return
	}

	// 保存交易记录
	saleRecord, err := generateRecord(req)
	if err != nil {
		logger.Warn("request parameter error->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}
	recordBean := &record.TxnRecord{
		BaseModel: models.BaseModel{
			Db:  models.DB(),
			Ctx: ctx,
		},
	}
	err = recordBean.Create(saleRecord)
	if err != nil {
		logger.Warn("create record error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	// 执行交易
	saleResp, errorCode := saleImp.Sale(ctx, &acquirer_impl.SaleRequest{
		TxqReq:       req,
		AgencyInfo:   agencyInfo,
		MerchantInfo: merchantInfo,
	})

	switch errorCode {
	case conf.Success: // success 逻辑写后面

	case conf.Reversal: // 需要冲正
		if err = saleRecord.UpdateStatus(record.NeedReversal); err != nil {
			logger.Error("update to reversal fail->", err.Error())
		}
		modules.BaseError(ctx, errorCode)
		return
	default:
		if err = saleRecord.UpdateStatus(record.Fail); err != nil {
			logger.Error("update to fail status fail->", err.Error())
		}
		modules.BaseError(ctx, errorCode)
		return
	}

	// Success

}

func generateRecord(req *api_define.TxnReq) (*record.TxnRecord, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, errors.New("can't parse amount->" + req.Amount + "," + err.Error())
	}

	ret := &record.TxnRecord{
		MerchantID:          req.MerchantID,
		TotalAmount:         amount,
		Amount:              amount,
		Currency:            req.Currency,
		MerchantAccountID:   req.PaymentProcessRule.MerchantID,
		PaymentMethod:       req.RealPaymentMethod,
		PaymentEntryType:    req.RealEntryType,
		PaymentType:         req.TxnType,
		PartnerUUID:         req.Uuid,
		Status:              record.Init,
		PaymentFromName:     req.FromName,
		PaymentFromIP:       req.FromIp,
		PaymentFromDeviceID: req.DeviceID,
		InvoiceNum:          "",
		CashierID:           req.CashierID,
	}

	if req.CreditCardBean != nil && req.CreditCardBean.CardNumber != "" {
		ret.ConsumerIdentify = req.CreditCardBean.CardNumber
	}

	return ret, nil
}

func fetchMerchantAgencyInfo(ctx *gin.Context, merchantID uint) (*merchant.Merchant, *agency.Agency, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)
	var err error

	merchantBean := &merchant.Merchant{
		BaseModel: models.BaseModel{
			Db:  models.DB(),
			Ctx: ctx,
		},
	}
	merchantBean, err = merchantBean.Get(merchantID)
	if err != nil {
		logger.Warn("merchant fetch fail->", err.Error())
		return nil, nil, conf.DBError
	}

	if merchantBean == nil {
		logger.Warn("can't find merchant")
		return nil, nil, conf.ProcessRuleSettingError
	}

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
		return nil, nil, conf.DBError
	}

	if agencyBean == nil {
		logger.Warn("can't find agency")
		return nil, nil, conf.ProcessRuleSettingError
	}

	return merchantBean, agencyBean, conf.Success
}
