package appindevice

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	tms2 "tpayment/modules/tms"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func UpdateHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(tms.AppInDevice)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 查询是否已经存在的账号
	bean, err := tms.GetAppInDeviceByID(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Info("GetDeviceByID sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if bean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	// 获取设备标识，查看是否有权限
	deviceBean, err := tms.GetDeviceByID(models.DB(), ctx, bean.ExternalId)
	if err != nil {
		logger.Error("GetDeviceByID fail->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if deviceBean == nil {
		logger.Error("device not found->", bean.ExternalId)
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}
	if err := tms2.CheckPermission(ctx, deviceBean); err != nil {
		logger.Error("CheckPermission fail->", err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	// 生成新账号
	req.ExternalId = 0
	req.ExternalIdType = ""
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
