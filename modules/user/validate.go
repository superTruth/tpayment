package user

import (
	"tpayment/conf"
	"tpayment/models/account"
	"tpayment/modules"

	"github.com/gin-gonic/gin"
)

// 登录
func ValidateHandle(ctx *gin.Context) {
	var userBean *account.UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*account.UserBean)
	} else {
		modules.BaseSuccess(ctx, conf.UnknownError)
		return
	}

	// 拼接response数据
	ret := &LoginResponse{
		Role:  userBean.Role,
		Name:  userBean.Name,
		Email: userBean.Email,
	}

	modules.BaseSuccess(ctx, ret)
}
