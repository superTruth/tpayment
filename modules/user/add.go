package user

import (
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"
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

	// 查询是否已经存在的账号
	user, err := account.GetUserByEmail(req.Email)
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
	account := &account.UserBean{
		Email: req.Email,
		Pwd:   req.Pwd,
		Name:  req.Name,
		Role:  req.Role,
	}

	err = models.CreateBaseRecord(account)

	if err != nil {
		logger.Info("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
