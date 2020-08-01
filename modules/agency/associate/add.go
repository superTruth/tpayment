package associate

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

func AddHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(agency.UserAgencyAssociate)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	if req.AgencyId == 0 || req.UserId == 0 {
		logger.Warn("ParameterError")
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 查询是否存在这2个ID
	userBean, err := account.GetUserById(models.DB(), ctx, req.UserId)
	if err != nil {
		logger.Info("GetUserById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if userBean == nil {
		logger.Warn("User Not Exist")
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	agencyBean, err := agency.GetAgencyById(models.DB(), ctx, req.AgencyId)
	if err != nil {
		logger.Info("GetAssociateById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if agencyBean == nil {
		logger.Warn("Agency Not Exist")
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
	if associateBean != nil {
		logger.Warn(conf.UserCanOnlyInOneAgency.String())
		modules.BaseError(ctx, conf.UserCanOnlyInOneAgency)
		return errors.New(conf.UserCanOnlyInOneAgency.String())
	}

	// 如果是系统管理员，则无法关联
	if userBean.Role != string(conf.RoleUser) { // 只有普通用户可以关联
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
