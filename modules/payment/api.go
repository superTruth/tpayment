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
	var err error

	req := new(api_define.TxnReq)

	err = utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	resp, errorCode := saleHandle(ctx, req)

	if errorCode != conf.Success {
		modules.BaseError(ctx, errorCode)
		return
	}
	modules.BaseSuccess(ctx, resp)
}

func VoidHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)
	var err error

	req := new(api_define.TxnReq)

	err = utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	resp, errorCode := voidHandle(ctx, req)

	if errorCode != conf.Success {
		modules.BaseError(ctx, errorCode)
		return
	}
	modules.BaseSuccess(ctx, resp)
}

func RefundHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)
	var err error

	req := new(api_define.TxnReq)

	err = utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	resp, errorCode := refundHandle(ctx, req)

	if errorCode != conf.Success {
		modules.BaseError(ctx, errorCode)
		return
	}
	modules.BaseSuccess(ctx, resp)
}
