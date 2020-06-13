package user

import (
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models/account"
	"tpayment/modules"
)

// 登录
func ValidateHandle(ctx echo.Context) error {
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)

	// 拼接response数据
	ret := &LoginResponse{
		Role:  userBean.Role,
		Name:  userBean.Name,
		Email: userBean.Email,
	}

	modules.BaseSuccess(ctx, ret)

	return nil
}
