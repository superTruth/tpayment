package agency

import (
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"
)

func UpdateHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(agency.Agency)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 查询是否已经存在的账号
	merchantBean, err := agency.GetAgencyById(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Info("GetAgencyById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if merchantBean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}

	// 生成新账号
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
