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

func AddHandle(ctx *gin.Context) {
	logger := tlog.GetGoroutineLogger()

	logger.Info("In App AddHandle")

	req := new(tms.AppInDevice)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}
	if req.ExternalId == 0 {
		logger.Warn("req.ExternalId == 0")
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 获取设备标识，查看是否有权限
	bean, err := tms.GetDeviceByID(req.ExternalId)
	if err != nil {
		logger.Error("GetDeviceByID fail->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if bean == nil {
		logger.Error("device not found->", req.ExternalId)
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}
	if err := tms2.CheckPermission(ctx, bean); err != nil {
		logger.Error("CheckPermission fail->", err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	errorCode := SmartAddAppInDevice(ctx, bean, req)
	if errorCode != conf.Success {
		modules.BaseError(ctx, errorCode)
		return
	}

	modules.BaseSuccess(ctx, nil)
}

// 智能添加app到设备
func SmartAddAppInDevice(ctx *gin.Context, device *tms.DeviceInfo, app *tms.AppInDevice) conf.ResultCode {
	logger := tlog.GetGoroutineLogger()

	// 查找是否已经存在这个app，如果存在，就更新当前规则，如果不存在，再创建新的记录
	bean, err := tms.FindAppInDevice(ctx, device.ID, app)
	if err != nil {
		logger.Error("FindAppInDevice fail->", err.Error())
		return conf.DBError
	}

	// 本来就不存在这个app的情况，就直接新增
	if bean == nil {
		logger.Info("app not exist->", app.AppFileId)
		bean := &tms.AppInDevice{
			ExternalId:     device.ID,
			ExternalIdType: tms.AppInDeviceExternalIdTypeDevice,
			Status:         app.Status,
			AppID:          app.AppID,
			AppFileId:      app.AppFileId,
		}
		if err = models.CreateBaseRecord(bean); err != nil {
			logger.Error("Create fail->", err.Error())
			return conf.DBError
		}
		return conf.Success
	}

	// app存在, 并且和原始配置一模一样，则不需要操作
	if (app.AppFileId == bean.AppFileId) && (app.Status == bean.Status) {
		logger.Info("app not change->", app.AppFileId)
		return conf.Success
	}

	// app需要update的情况
	bean.Status = app.Status
	bean.AppFileId = app.AppFileId
	bean.AppID = app.AppID

	if err := models.UpdateBaseRecord(bean); err != nil {
		logger.Error("Update fail->", err.Error)
		return conf.DBError
	}

	return conf.Success
}
