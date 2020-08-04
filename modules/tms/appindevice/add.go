package appindevice

import (
	"errors"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	tms2 "tpayment/modules/tms"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func AddHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(tms.AppInDevice)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil || req.ExternalId == 0 {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 获取设备标识，查看是否有权限
	bean, err := tms.GetDeviceByID(models.DB(), ctx, req.ExternalId)
	if err != nil {
		logger.Error("GetDeviceByID fail->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if bean == nil {
		logger.Error("device not found->", req.ExternalId)
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}
	if err := tms2.CheckPermission(ctx, bean); err != nil {
		logger.Error("CheckPermission fail->", err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return err
	}

	errorCode := SmartAddAppInDevice(ctx, req)
	if errorCode != conf.SUCCESS {
		modules.BaseError(ctx, errorCode)
		return errors.New(errorCode.String())
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}

// 智能添加app到设备
func SmartAddAppInDevice(ctx echo.Context, app *tms.AppInDevice) conf.ResultCode {
	logger := tlog.GetLogger(ctx)

	// 查找是否已经存在这个app，如果存在，就更新当前规则，如果不存在，再创建新的记录
	bean, err := tms.FindAppInDevice(models.DB(), ctx, app)
	if err != nil {
		logger.Error("FindAppInDevice fail->", err.Error())
		return conf.DBError
	}

	// 本来就不存在这个app的情况，就直接新增
	if bean == nil {
		logger.Info("app not exist->", app.AppFileId)
		bean := &tms.AppInDevice{
			ExternalId:     app.ExternalId,
			ExternalIdType: tms.AppInDeviceExternalIdTypeDevice,
			Status:         app.Status,
			AppID:          app.AppID,
			AppFileId:      app.AppFileId,
		}
		if err = models.CreateBaseRecord(bean); err != nil {
			logger.Error("Create fail->", err.Error())
			return conf.DBError
		}
		return conf.SUCCESS
	}

	// app存在, 并且和原始配置一模一样，则不需要操作
	if (app.AppFileId == bean.AppFileId) && (app.Status == bean.Status) {
		logger.Info("app not change->", app.AppFileId)
		return conf.SUCCESS
	}

	// app需要update的情况
	bean.Status = app.Status
	bean.AppFileId = app.AppFileId
	bean.AppID = app.AppID

	if err := models.UpdateBaseRecord(bean); err != nil {
		logger.Error("Update fail->", err.Error)
		return conf.DBError
	}

	return conf.SUCCESS
}
