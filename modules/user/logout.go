package user

import (
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/tlog"
)

func LogoutHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	token := ctx.Request().Header[conf.HeaderTagToken][0]

	tokenBean, err := account.GetTokenBeanByToken(token)

	if err != nil {
		logger.Error("GetTokenBeanByToken sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	if tokenBean != nil {
		_ = models.DeleteBaseRecord(tokenBean)
	}

	//
	ret := &modules.BaseResponse{
		ErrorCode:    conf.SUCCESS,
	}

	modules.BaseSuccess(ctx, ret)

	return nil
}
