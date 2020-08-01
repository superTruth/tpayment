package acquirer

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func AddHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(agency.Acquirer)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 管理员必须要传入agency id
	var agencyId uint
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
	if userBean.Role == string(conf.RoleAdmin) {
		if req.AgencyId == 0 {
			logger.Warn("Admin user must contain agency id->")
			modules.BaseError(ctx, conf.ParameterError)
			return err
		}
		agencyId = req.AgencyId
	} else {
		agencys := ctx.Get(conf.ContextTagAgency).([]*agency.Agency)
		agencyId = agencys[0].ID
		req.AgencyId = agencyId
	}

	agencyBean, err := agency.GetAgencyById(models.DB(), ctx, agencyId)
	if err != nil {
		logger.Error("GetAgencyById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	if agencyBean == nil {
		logger.Info("the agency is not found->", agencyId)
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
