package merchantdevicepayment

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/merchant"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func DeleteHandle(ctx *gin.Context) {
	logger := tlog.GetGoroutineLogger()

	req := new(modules.BaseIDRequest)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	deviceBean, err := merchant.GetPaymentSettingInDeviceById(req.ID)
	if err != nil {
		logger.Error("GetPaymentSettingInDeviceById fail->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	if deviceBean == nil {
		logger.Error(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	err = models.DeleteBaseRecord(deviceBean)

	if err != nil {
		logger.Info("DeleteBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
