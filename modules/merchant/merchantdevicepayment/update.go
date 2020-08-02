package merchantdevicepayment

import (
	"errors"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/merchant"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func UpdateHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(merchant.PaymentSettingInDevice)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	deviceBean, err := merchant.GetPaymentSettingInDeviceById(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Error("GetPaymentSettingInDeviceById fail->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	if deviceBean == nil {
		logger.Error(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return errors.New(conf.RecordNotFund.String())
	}

	// TODO 权限判断
	req.MerchantDeviceId = 0
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
