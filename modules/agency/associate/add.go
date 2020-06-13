package associate

import (
	"errors"
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"
)

func AddHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(agency.UserAgencyAssociate)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 一个用户只可以关联一个agency
	associateBean, err := agency.GetAssociateByUserId(models.DB(), ctx, req.UserId)
	if err != nil {
		logger.Info("GetAssociateByUserId sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if associateBean == nil {
		logger.Warn(conf.UserCanOnlyInOneAgency)
		modules.BaseError(ctx, conf.UserCanOnlyInOneAgency)
		return errors.New(conf.UserCanOnlyInOneAgency.String())
	}

	// 如果是系统管理员，则无法关联
	userBean, err := account.GetUserById(req.UserId)
	if err != nil {
		logger.Info("GetUserById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if userBean == nil { // 账号无法找到
		logger.Info("cant find the user")
		modules.BaseError(ctx, conf.RecordNotFund)
		return errors.New(conf.RecordNotFund.String())
	}
	if userBean.Role == string(conf.RoleAdmin) { // 系统管理员已经拥有所有权限，所以不运行添加关联
		logger.Info(conf.AdminCantAssociate.String())
		modules.BaseError(ctx, conf.AdminCantAssociate)
		return errors.New(conf.AdminCantAssociate.String())
	}

	err = models.CreateBaseRecord(req)
	if err != nil {
		logger.Info("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
