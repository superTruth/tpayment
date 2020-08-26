package user

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func LogoutHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	token := ctx.Request.Header[conf.HeaderTagToken][0]

	tokenBean, err := account.GetTokenBeanByToken(models.DB(), ctx, token)

	if err != nil {
		logger.Error("GetTokenBeanByToken sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	if tokenBean != nil {
		_ = models.DeleteBaseRecord(tokenBean)
	}

	modules.BaseSuccess(ctx, nil)
}
