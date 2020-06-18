package appindevice

import (
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/models/merchant"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"
)

// TODO 未完成
func AddHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(merchant.Merchant)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// TODO 未做判断：当前用户可能没有此机构权限
	agencyBean, err := agency.GetAgencyById(models.DB(), ctx, req.AgencyId)
	if err != nil {
		logger.Error("GetAgencyById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	if agencyBean == nil {
		logger.Info("GetAgencyById sql error->", err.Error())
		modules.BaseError(ctx, conf.RecordNotFund)
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
