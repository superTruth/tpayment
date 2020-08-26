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

func UpdateHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(merchant.PaymentSettingInDevice)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	deviceBean, err := merchant.GetPaymentSettingInDeviceById(models.DB(), ctx, req.ID)
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

	req.MerchantDeviceId = 0
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
