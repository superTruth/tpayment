package void

import (
	"strconv"
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

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

const saleMaxExpTime = time.Minute * 5

func Handle(ctx *gin.Context) {
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

	// 查找原始记录
	recordBean := record.TxnRecord{
		BaseModel: models.BaseModel{
			Db:  models.DB(),
			Ctx: ctx,
		},
	}
	req.OrgRecord, err = recordBean.GetByID(req.OriginTxnID)
	if err != nil {
		logger.Warn("GetByID "+strconv.Itoa(int(req.OriginTxnID))+" fail->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if req.OrgRecord == nil {
		logger.Warn("can't find the record " + strconv.Itoa(int(req.OriginTxnID)))
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}
	req.OrgRecord.BaseModel = recordBean.BaseModel

	// 判断是否可以void
	if req.OrgRecord.AcquirerSettlementAt != nil { // 被结算过的交易不能void
		logger.Warn("the record was settled")
		modules.BaseError(ctx, conf.Settled)
		return
	}

	if req.OrgRecord.VoidAt != nil { // 已经被void过
		logger.Warn("the record was voided")
		modules.BaseError(ctx, conf.Voided)
		return
	}

	// 获取merchant account, acquirer
	errorCode = fetchMerchantAccount(ctx, req)
	if errorCode != conf.Success {
		logger.Warn("fetchMerchantAccount fail->", errorCode.String())
		modules.BaseError(ctx, errorCode)
		return
	}
	if req.PaymentProcessRule.MerchantAccount.Terminal != nil { // 如果有TID的情况，需要锁定TID
		errorCode = req.PaymentProcessRule.MerchantAccount.Terminal.Lock(saleMaxExpTime)
		if errorCode != conf.Success {
			modules.BaseError(ctx, errorCode)
			return
		}
		defer req.PaymentProcessRule.MerchantAccount.Terminal.UnLock()
	}

	// 获取void交易对象
	acquirerImpl, ok := acquirer_impl.AcquirerImpls[req.PaymentProcessRule.MerchantAccount.Acquirer.Name]
	if !ok {
		logger.Warn("can't find acquirer impl->", req.PaymentProcessRule.MerchantAccount.Acquirer.Name)
		modules.BaseError(ctx, conf.UnknownError)
		return
	}
	voidImp, ok := acquirerImpl.(acquirer_impl.IVoid)
	if !ok {
		logger.Warn("the acquirer not support sale->", req.PaymentProcessRule.MerchantAccount.Acquirer.Name)
		modules.BaseError(ctx, conf.UnknownError)
		return
	}

	// 保存交易记录
	req.TxnRecord, err = buildRecord(req)
	if err != nil {
		logger.Warn("request parameter error->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}
	req.TxnRecord.BaseModel = models.BaseModel{
		Db:  models.DB(),
		Ctx: ctx,
	}
	err = req.TxnRecord.Create(req.TxnRecord)
	if err != nil {
		logger.Warn("create record error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	// 创建response数据
	resp := preBuildResp(req)

	// 执行交易
	saleResp, errorCode := voidImp.Void(ctx, &acquirer_impl.SaleRequest{
		TxqReq:       req,
		AgencyInfo:   agencyInfo,
		MerchantInfo: merchantInfo,
	})

	switch errorCode {
	case conf.Success: // success 逻辑写后面

	case conf.Reversal: // 需要冲正
		if err = req.TxnRecord.UpdateStatus(record.NeedReversal); err != nil {
			logger.Error("update to reversal fail->", err.Error())
		}
		modules.BaseError(ctx, errorCode)
		return
	default:
		if err = req.TxnRecord.UpdateStatus(record.Fail); err != nil {
			logger.Error("update to fail status fail->", err.Error())
		}
		modules.BaseError(ctx, errorCode)
		return
	}

	// Success，合并response
	mergeAcquirerResponse(resp, saleResp)
	mergeResponseToRecord(req, saleResp)

	// 交易成功，还需要同时更新原始交易，作为事物
	if req.TxnRecord.Status == record.Success {
		t := time.Now()
		req.OrgRecord.VoidAt = &t
		err = models.DB().Transaction(func(tx *gorm.DB) error {
			// 原始记录
			req.OrgRecord.BaseModel = models.BaseModel{
				Db:  &models.MyDB{tx},
				Ctx: ctx,
			}

			err = req.OrgRecord.UpdateVoidStatus()
			if err != nil {
				logger.Error("UpdateVoidStatus fail->", err.Error())
				return err
			}

			// 新记录
			req.TxnRecord.BaseModel = models.BaseModel{
				Db:  &models.MyDB{tx},
				Ctx: ctx,
			}

			err = req.TxnRecord.UpdateTxnResult()
			if err != nil {
				logger.Error("UpdateTxnResult fail->", err.Error())
				return err
			}

			return nil
		})

		if err != nil {
			logger.Error("update success result fail->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return
		}
	} else {
		if err = req.TxnRecord.UpdateTxnResult(); err != nil {
			logger.Error("UpdateTxnResult fail->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return
		}
	}

	modules.BaseSuccess(ctx, saleResp)
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
