package middleware

import (
	"strings"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/modules/user"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func AuthHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	if ctx.Request.RequestURI == conf.UrlAccountLogin ||
		ctx.Request.RequestURI == conf.UrlAccountRegister ||
		strings.Contains(ctx.Request.RequestURI, "/payment/account/active") { // 唯一的登录功能不需要token
		ctx.Next()
		return
	}

	// 判断token
	tokens := ctx.Request.Header[conf.HeaderTagToken]
	if len(tokens) == 0 {
		logger.Info("authHandle error->", conf.NeedTokenInHeader.String())
		modules.BaseError(ctx, conf.NeedTokenInHeader)
		ctx.Next()
		return
	}

	// 验证token的有效性
	userBean, _, err := user.Auth(ctx, tokens[0])
	if err != nil { // 数据库出错
		logger.Error("Auth db error2->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		ctx.Next()
		return
	}
	if userBean == nil { // token验证失败
		logger.Error("authHandle token not exist")
		modules.BaseError(ctx, conf.TokenInvalid)
		ctx.Next()
		return
	}

	ctx.Set(conf.ContextTagUser, userBean)

	// 查看是否是机构管理员
	if userBean.Role == string(conf.RoleUser) { // 机器人和系统管理员不需要验证
		_, agencyBean, err := agency.QueryAgencyRecord(models.DB(), ctx, 0, 1000, nil)
		if err != nil {
			logger.Error("QueryAgencyRecord db error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			ctx.Next()
			return
		}
		ctx.Set(conf.ContextTagAgency, agencyBean)
	} else {
		ctx.Set(conf.ContextTagAgency, make([]*agency.Agency, 0))
	}

	ctx.Next()
}
