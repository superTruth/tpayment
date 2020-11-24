package payment

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func SaleHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(api_define.TxnReq)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 预处理交易数据，分析出真实卡号，用卡方式
	errorCode := preHandleRequest(ctx, req)
	if errorCode != conf.SUCCESS {
		logger.Warn("preHandleRequest fail->", errorCode.String())
		modules.BaseError(ctx, errorCode)
		return
	}

	// 提取payment processing rule
	req.PaymentProcessRule, errorCode = matchProcessRule(ctx, req)
	if errorCode != conf.SUCCESS {
		logger.Warn("matchProcessRule fail->", errorCode.String())
		modules.BaseError(ctx, errorCode)
		return
	}

	// 获取merchant account, acquirer
	errorCode = fetchMerchantAccount(ctx, req)
	if errorCode != conf.SUCCESS {
		logger.Warn("fetchMerchantAccount fail->", errorCode.String())
		modules.BaseError(ctx, errorCode)
		return
	}

	// 保存交易记录

	// 执行交易

}
