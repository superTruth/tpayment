package user

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

	req := new(AddUserRequest)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 机构管理员
	agencyId := uint(0)
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
	agencys := ctx.Get(conf.ContextTagAgency).([]*agency.Agency)
	if userBean.Role != string(conf.RoleAdmin) { // 管理员，不需要过滤机构
		agencyId = agencys[0].ID
	}

	// 查询是否已经存在的账号
	user, err := account.GetUserByEmail(models.DB(), ctx, req.Email)
	if err != nil {
		logger.Info("GetUserByEmail sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if user != nil {
		logger.Warn(conf.RecordAlreadyExist.String())
		modules.BaseError(ctx, conf.RecordAlreadyExist)
		return err
	}

	// 生成新账号
	bean := &account.UserBean{
		AgencyId: agencyId,
		Email:    req.Email,
		Pwd:      req.Pwd,
		Name:     req.Name,
		Role:     req.Role,
		Active:   true,
	}

	err = models.CreateBaseRecord(bean)

	if err != nil {
		logger.Info("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
