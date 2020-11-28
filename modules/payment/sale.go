package payment

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

	"github.com/gin-gonic/gin"
)

const saleMaxExpTime = time.Minute * 5

func saleHandle(ctx *gin.Context, req *api_define.TxnReq) (*api_define.TxnResp, conf.ResultCode) {
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
		errorCode = req.PaymentProcessRule.MerchantAccount.Terminal.Lock(saleMaxExpTime)
		if errorCode != conf.Success {
			return resp, errorCode
		}
		defer func() {
			req.PaymentProcessRule.MerchantAccount.Terminal.UnLock()
		}()
	}

	// 获取sale交易对象
	acquirerImpl, ok := factory.AcquirerImpls[req.PaymentProcessRule.MerchantAccount.Acquirer.Name]
	if !ok {
		logger.Warn("can't find acquirer impl->", req.PaymentProcessRule.MerchantAccount.Acquirer.Name)
		return resp, conf.UnknownError
	}
	saleImp, ok := acquirerImpl.(acquirer_impl.ISale)
	if !ok {
		logger.Warn("the acquirer not support sale->", req.PaymentProcessRule.MerchantAccount.Acquirer.Name)
		return resp, conf.UnknownError
	}

	logger.Info("save.....", (req.TxnRecord == nil))
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

	logger.Info("mergeRespAfterPreHandle.....")
	// 再次合并数据到返回结果
	mergeRespAfterPreHandle(resp, req)

	logger.Info("执行交易.....")
	// 执行交易
	saleResp, errorCode := saleImp.Sale(ctx, &acquirer_impl.SaleRequest{
		TxqReq: req,
	})

	logger.Info("txn result->", errorCode)
	switch errorCode {
	case conf.Success: // success 逻辑写后面

	case conf.Reversal: // 需要冲正
		if err = req.TxnRecord.UpdateStatus(record.NeedReversal); err != nil {
			logger.Error("update to reversal fail->", err.Error())
		}
		modules.BaseError(ctx, errorCode)
		return resp, errorCode
	default:
		if err = req.TxnRecord.UpdateStatus(record.Fail); err != nil {
			logger.Error("update to fail status fail->", err.Error())
		}
		return resp, errorCode
	}

	// Success，合并response
	mergeAcquirerResponse(resp, saleResp)
	mergeResponseToRecord(req.TxnRecord, saleResp)
	if err = req.TxnRecord.UpdateTxnResult(); err != nil {
		logger.Error("UpdateTxnResult fail->", err.Error())
		return resp, conf.DBError
	}
	return resp, conf.Success
}
