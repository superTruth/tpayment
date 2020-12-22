package pay_online

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/internal/acquirer_impl/factory"
	"tpayment/models"
	"tpayment/models/payment/record"
	"tpayment/modules"
	"tpayment/pkg/id"
	"tpayment/pkg/tlog"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

func transferWithOnlinePaymentHandle(ctx *gin.Context, req *api_define.TxnReq) (*api_define.TxnResp, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)
	var err error

	err = api_define.Validate(ctx, req)
	if err != nil {
		logger.Warn("validate request body error->", err.Error())
		return nil, conf.ParameterError
	}

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
		logger.Info("lock tid->", req.PaymentProcessRule.MerchantAccount.Terminal.TID)
		errorCode = req.PaymentProcessRule.MerchantAccount.Terminal.Lock(saleMaxExpTime)
		if errorCode != conf.Success {
			return resp, errorCode
		}
		defer func() {
			logger.Info("unlock tid->", req.PaymentProcessRule.MerchantAccount.Terminal.TID)
			req.PaymentProcessRule.MerchantAccount.Terminal.UnLock()
		}()
		req.CreditCardBean.TraceNum = req.PaymentProcessRule.MerchantAccount.Terminal.TraceNum
		req.CreditCardBean.BatchNum = req.PaymentProcessRule.MerchantAccount.Terminal.BatchNum
		req.TxnRecord.AcquirerTraceNum = req.CreditCardBean.TraceNum
		req.TxnRecord.AcquirerBatchNum = req.CreditCardBean.BatchNum

		// trace No自增
		err = req.PaymentProcessRule.MerchantAccount.Terminal.IncTraceNum()
		if err != nil {
			logger.Error("req.PaymentProcessRule.MerchantAccount.Terminal.IncTraceNum fail->", err.Error())
			return resp, conf.DBError
		}
	}

	// 获取sale交易对象
	acquirerImpl, ok := factory.AcquirerImpls[req.PaymentProcessRule.MerchantAccount.Acquirer.ImplName]
	if !ok {
		logger.Warn("can't find acquirer impl->", req.PaymentProcessRule.MerchantAccount.Acquirer.Name)
		return resp, conf.UnknownError
	}
	saleImp, ok := acquirerImpl.(acquirer_impl.ISale)
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

	// 修正交易
	err = models.DB().Transaction(func(tx *gorm.DB) error {
		req.TxnRecord.Db = &models.MyDB{DB: tx}
		if err = req.TxnRecord.UpdateTxnResult(); err != nil {
			logger.Error("UpdateTxnResult fail->", err.Error())
			return err
		}

		// 创建一条transfer记录
		transferRecord := &record.TxnRecord{
			BaseModel: models.BaseModel{
				ID: id.New(),
				Db: req.TxnRecord.Db,
			},
			MerchantID:  req.TxnRecord.MerchantID,
			PaymentType: conf.Transfer,
			OrgTxnID:    req.TxnRecord.ID,
			Status:      record.Init,
		}

		if err = transferRecord.Create(transferRecord); err != nil {
			logger.Error("Create fail->", err.Error())
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error("UpdateTxnResult fail->", err.Error())
		return resp, conf.DBError
	}
	return resp, conf.Success
}
