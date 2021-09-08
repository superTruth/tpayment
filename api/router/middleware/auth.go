package middleware

import (
	"strings"
	"tpayment/conf"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/modules/user"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func AuthHandle(ctx *gin.Context) {
	logger := tlog.GetGoroutineLogger()

	if ctx.Request.RequestURI == conf.UrlAccountLogin ||
		ctx.Request.RequestURI == conf.UrlAccountRegister ||
		strings.Contains(ctx.Request.RequestURI, "/payment/account/active") { // 唯一的登录功能不需要token
		ctx.Next()
		return
	}

	// 判断token
	if len(ctx.Request.Header[conf.HeaderTagToken]) == 0 &&
		(len(ctx.Request.Header[conf.HeaderTagAccessKey]) == 0 ||
			len(ctx.Request.Header[conf.HeaderTagAccessHash]) == 0) {
		logger.Info("authHandle error->", conf.NeedTokenInHeader.String())
		modules.BaseError(ctx, conf.NeedTokenInHeader)
		ctx.Abort()
		return
	}

	// 验证登录权限
	userBean, _, err := user.Auth(ctx)
	if err != nil { // 数据库出错
		logger.Error("Auth db error2->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		ctx.Abort()
		return
	}
	if userBean == nil { // token验证失败
		logger.Error("authHandle token not exist")
		modules.BaseError(ctx, conf.TokenInvalid)
		ctx.Abort()
		return
	}

	ctx.Set(conf.ContextTagUser, userBean)

	// 查看是否是机构管理员
	if userBean.Role == string(conf.RoleUser) { // 机器人和系统管理员不需要验证
		_, agencyBean, err := agency.QueryAgencyRecord(ctx, 0, 1000, nil)
		if err != nil {
			logger.Error("QueryAgencyRecord db error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			ctx.Abort()
			return
		}
		ctx.Set(conf.ContextTagAgency, agencyBean)
	} else {
		ctx.Set(conf.ContextTagAgency, make([]*agency.Agency, 0))
	}

	ctx.Next()
}
