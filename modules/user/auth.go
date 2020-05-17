package user

import (
	"github.com/labstack/echo"
	"tpayment/models/account"
	"tpayment/pkg/tlog"
)

func Auth(ctx echo.Context, token string) (*account.UserBean, *account.AppIdBean, error) {
	logger := tlog.GetLogger(ctx)
	// 创建 或者 更新  token记录
	tokenBean, err := account.GetTokenBeanByToken(token)
	if err != nil {
		logger.Warn("GetTokenBeanByToken fail->", err.Error())
		return nil, nil, err
	}

	//
	accountBean, err := account.GetUserById(tokenBean.UserId)
	if err != nil {
		logger.Warn("GetUserById fail->", err.Error())
		return nil, nil, err
	}

	appBean, err := account.GetAppIdByID(tokenBean.AppId)
	if err != nil {
		logger.Warn("GetUserById fail->", err.Error())
		return nil, nil, err
	}

	return accountBean, appBean, nil
}
