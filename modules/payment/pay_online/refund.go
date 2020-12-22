package pay_online

import (
	"time"
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/internal/acquirer_impl/factory"
	"tpayment/models"
	"tpayment/models/payment/record"
	"tpayment/modules"
	"tpayment/pkg/tlog"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

const refundMaxExpTime = time.Minute * 5

func refundHandle(ctx *gin.Context, req *api_define.TxnReq) (*api_define.TxnResp, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)
	var err error

	// 创建response数据
	resp := preBuildResp(req)

	// 预处理请求数据，解析卡数据
	errorCode := preHandleRequest(ctx, req)
	if errorCode != conf.Success {
		logger.Warn("preHandleRequest fail->", errorCode.String())
		return resp, errorCode
	}

	// 锁定TID
	if req.PaymentProcessRule.MerchantAccount.Terminal != nil { // 如果有TID的情况，需要锁定TID
		errorCode = req.PaymentProcessRule.MerchantAccount.Terminal.Lock(refundMaxExpTime)
		if errorCode != conf.Success {
			return resp, errorCode
		}
		defer req.PaymentProcessRule.MerchantAccount.Terminal.UnLock()
	}

	// 获取void交易对象
	acquirerImpl, ok := factory.AcquirerImpls[req.PaymentProcessRule.MerchantAccount.Acquirer.ImplName]
	if !ok {
		logger.Warn("can't find acquirer impl->", req.PaymentProcessRule.MerchantAccount.Acquirer.Name)
		return resp, conf.UnknownError
	}
	refundImp, ok := acquirerImpl.(acquirer_impl.IRefund)
	if !ok {
		logger.Warn("the acquirer not support sale->", req.PaymentProcessRule.MerchantAccount.Acquirer.Name)
		return resp, conf.UnknownError
	}

	// 保存交易记录
	req.TxnRecord.BaseModel = models.BaseModel{
		Db:  models.DB(),
		Ctx: ctx,
	}
	err = req.TxnRecord.Create(req.TxnRecord)
	if err != nil {
		logger.Warn("create record error->", err.Error())
		return resp, conf.DBError
	}

	// 执行交易
	saleResp, errorCode := refundImp.Refund(ctx, &acquirer_impl.SaleRequest{
		TxqReq: req,
	})

	switch errorCode {
	case conf.Success: // success 逻辑写后面

	case conf.Reversal: // 需要冲正
		req.TxnRecord.Status = record.NeedReversal
		req.TxnRecord.ErrorCode = string(errorCode)
		req.TxnRecord.ErrorDes = errorCode.String()
		if err = req.TxnRecord.UpdateStatus(); err != nil {
			logger.Error("update to reversal fail->", err.Error())
		}
		modules.BaseError(ctx, errorCode)
		return resp, errorCode
	default:
		req.TxnRecord.Status = record.Fail
		req.TxnRecord.ErrorCode = string(errorCode)
		req.TxnRecord.ErrorDes = errorCode.String()
		if err = req.TxnRecord.UpdateStatus(); err != nil {
			logger.Error("update to fail status fail->", err.Error())
		}
		return resp, errorCode
	}

	// Success，合并response
	mergeAcquirerResponse(resp, saleResp)
	mergeResponseToRecord(req, saleResp)

	if req.TxnRecord.Status == record.Success {
		t := time.Now()
		req.OrgRecord.RefundAt = &t
		req.OrgRecord.RefundTimes++
		req.OrgRecord.TotalAmount = req.OrgRecord.TotalAmount.Sub(req.TxnRecord.Amount)
		err = models.DB().Transaction(func(tx *gorm.DB) error {
			// 原始记录
			req.OrgRecord.BaseModel = models.BaseModel{
				Db:  &models.MyDB{DB: tx},
				Ctx: ctx,
			}

			err = req.OrgRecord.UpdateRefundStatus()
			if err != nil {
				logger.Error("UpdateRefundStatus fail->", err.Error())
				return err
			}

			// 新记录
			req.TxnRecord.BaseModel = models.BaseModel{
				Db:  &models.MyDB{DB: tx},
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
			return resp, conf.DBError
		}
	} else {
		if err = req.TxnRecord.UpdateTxnResult(); err != nil {
			logger.Error("UpdateTxnResult fail->", err.Error())
			return resp, conf.DBError
		}
	}

	return resp, conf.Success
}
