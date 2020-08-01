package user

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/tlog"

	"github.com/labstack/echo"
)

func LogoutHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	token := ctx.Request().Header[conf.HeaderTagToken][0]

	tokenBean, err := account.GetTokenBeanByToken(models.DB(), ctx, token)

	if err != nil {
		logger.Error("GetTokenBeanByToken sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	if tokenBean != nil {
		_ = models.DeleteBaseRecord(tokenBean)
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
