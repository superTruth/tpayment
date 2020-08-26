package appinbatchupdate

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AddHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(tms.AppInDevice)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	req.ID = 0
	req.ExternalIdType = tms.AppInDeviceExternalIdTypeDevice
	err = models.CreateBaseRecord(req)
	if err != nil {
		logger.Error("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
