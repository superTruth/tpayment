package acquirer

import (
	"errors"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/models/agency"
	"tpayment/modules"
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

	// 查询是否已经存在的账号
	acquirer, err := agency.GetAcquirerById(req.ID)
	if err != nil {
		logger.Info("GetUserById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if acquirer == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}

	// 判断当前agency是否有权限
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
	if userBean.Role != string(conf.RoleAdmin) {
		agencys := ctx.Get(conf.ContextTagAgency).([]*agency.Agency)
		if agencys[0].ID != acquirer.AgencyId {
			logger.Warn("this acquirer is not belong to the agency")
			modules.BaseError(ctx, conf.NoPermission)
			return errors.New("this acquirer is not belong to the agency")
		}
	}

	err = models.DeleteBaseRecord(acquirer)

	if err != nil {
		logger.Info("DeleteBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
