package accesskey

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func DeleteHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(modules.BaseIDRequest)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 查询是否已经存在的账号
	keyBean, err := account.GetUserAccessKeyFromID(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Info("GetUserAccessKeyFromID sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if keyBean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	var userBean *account.UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*account.UserBean)
	} else {
		modules.BaseError(ctx, conf.UnknownError)
		return
	}

	// 查看是否是自己的账号，如果不是，不允许删除
	if userBean.ID != keyBean.UserId {
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	err = models.DeleteBaseRecord(keyBean)

	if err != nil {
		logger.Info("DeleteBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
