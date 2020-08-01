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

func UpdateHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(agency.Acquirer)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 查询是否已经存在的账号
	acquirerBean, err := agency.GetAcquirerById(req.ID)
	if err != nil {
		logger.Info("GetMerchantById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if acquirerBean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}

	// 判断当前agency是否有权限
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
	if userBean.Role != string(conf.RoleAdmin) {
		agencys := ctx.Get(conf.ContextTagAgency).([]*agency.Agency)
		if agencys[0].ID != acquirerBean.AgencyId {
			logger.Warn("this acquirer is not belong to the agency")
			modules.BaseError(ctx, conf.NoPermission)
			return errors.New("this acquirer is not belong to the agency")
		}
	}

	// 生成新账号
	req.AgencyId = 0 // 不允许更新agency id
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
