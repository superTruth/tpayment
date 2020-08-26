package uploadfile

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

	req := new(tms.UploadFile)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 现在找到对应的机器，看看机器所属的机构
	deviceBean, err := tms.GetDeviceBySn(models.DB(), ctx, req.DeviceSn)
	if err != nil {
		logger.Warn("GetDeviceBySn fail->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	if deviceBean == nil {
		logger.Warn("can't find the device->", req.DeviceSn)
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	req.AgencyId = deviceBean.AgencyId
	req.ID = 0
	err = models.CreateBaseRecord(req)

	if err != nil {
		logger.Error("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
