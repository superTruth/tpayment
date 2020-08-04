package appindevice

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	tms2 "tpayment/modules/tms"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func UpdateHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(tms.AppInDevice)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 查询是否已经存在的账号
	bean, err := tms.GetAppInDeviceByID(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Info("GetDeviceByID sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if bean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}

	// 获取设备标识，查看是否有权限
	deviceBean, err := tms.GetDeviceByID(models.DB(), ctx, bean.ExternalId)
	if err != nil {
		logger.Error("GetDeviceByID fail->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if deviceBean == nil {
		logger.Error("device not found->", bean.ExternalId)
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}
	if err := tms2.CheckPermission(ctx, deviceBean); err != nil {
		logger.Error("CheckPermission fail->", err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return err
	}

	// 生成新账号
	req.ExternalId = 0
	req.ExternalIdType = ""
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
