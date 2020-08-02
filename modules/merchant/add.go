package merchant

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/merchant"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func AddHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(merchant.Merchant)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 判断当前agency id
	req.AgencyId, err = modules.GetAgencyId(ctx, req.AgencyId)
	if err != nil {
		logger.Warn(err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	err = models.CreateBaseRecord(req)

	if err != nil {
		logger.Error("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
