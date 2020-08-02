package merchantdevice

import (
	"errors"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/merchant"
	"tpayment/modules"
	merchantModule "tpayment/modules/merchant"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func DeleteHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(modules.BaseIDRequest)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	deviceBean, err := merchant.GetDeviceInMerchantAssociateById(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Error("GetDeviceInMerchantAssociateById fail->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	if deviceBean == nil {
		logger.Error(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return errors.New(conf.RecordNotFund.String())
	}

	// 判断权限
	err = merchantModule.CheckPermission(ctx, deviceBean.MerchantId)
	if err != nil {
		logger.Warn(err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return err
	}

	bean := &merchant.DeviceInMerchant{
		BaseModel: models.BaseModel{
			ID: req.ID,
		},
	}
	err = models.DeleteBaseRecord(bean)

	if err != nil {
		logger.Info("DeleteBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
