package associate

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/merchant"
	"tpayment/modules"
	merchantModule "tpayment/modules/merchant"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AddHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(merchant.UserMerchantAssociate)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 判断权限
	err = merchantModule.CheckPermission(ctx, req.MerchantId)
	if err != nil {
		logger.Warn(err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	err = models.CreateBaseRecord(req)

	if err != nil {
		logger.Info("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
