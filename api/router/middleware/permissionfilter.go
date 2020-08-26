package middleware

import (
	"strings"
	"tpayment/conf"
	"tpayment/models/account"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func PermissionFilter(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	if ctx.Request.RequestURI == conf.UrlAccountLogin ||
		ctx.Request.RequestURI == conf.UrlAccountRegister ||
		strings.Contains(ctx.Request.RequestURI, "/payment/account/active") { // 唯一的登录功能不需要token
		ctx.Next()
		return
	}

	var userBean *account.UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*account.UserBean)
	} else {
		logger.Error("user is null")
		modules.BaseError(ctx, conf.UnknownError)
		ctx.Next()
		return
	}

	// 管理员直行
	if userBean.Role == string(conf.RoleAdmin) {
		ctx.Next()
		return
	}

	// 需要管理员权限的
	if ctx.Request.RequestURI == conf.UrlAgencyAdd || // 所有的机构操作只能管理员
		ctx.Request.RequestURI == conf.UrlAgencyUpdate ||
		ctx.Request.RequestURI == conf.UrlAgencyDelete ||
		strings.Contains(ctx.Request.RequestURI, "/payment/agency_device") ||
		strings.Contains(ctx.Request.RequestURI, "/payment/tms/model") {
		logger.Warn("no agency permission")
		modules.BaseError(ctx, conf.NoPermission)
		ctx.Next()
		return
	}

	// 需要管理员或者机构管理员的权限的
	var agencys []*agency.Agency
	agencysTmp, ok := ctx.Get(conf.ContextTagAgency)
	if ok {
		agencys = agencysTmp.([]*agency.Agency)
	}

	if ctx.Request.RequestURI == conf.UrlAccountAdd ||
		ctx.Request.RequestURI == conf.UrlAccountDelete ||
		ctx.Request.RequestURI == conf.UrlAccountQuery ||
		ctx.Request.RequestURI == conf.UrlMerchantAdd ||
		ctx.Request.RequestURI == conf.UrlMerchantUpdate ||
		ctx.Request.RequestURI == conf.UrlMerchantQuery ||
		strings.Contains(ctx.Request.RequestURI, "/payment/tms") {

		if userBean.Role == string(conf.RoleUser) { // 管理员，不需要过滤机构
			if len(agencys) == 0 {
				logger.Info("not admin or not agency")
				modules.BaseError(ctx, conf.NoPermission)
				ctx.Next()
				return
			}
		}
	}

	ctx.Next()
}
