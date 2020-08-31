package accesskey

import (
	"encoding/hex"
	"strings"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/algorithmutils"
	"tpayment/pkg/tlog"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func AddHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	var userBean *account.UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*account.UserBean)
	} else {
		logger.Error("user is null")
		modules.BaseError(ctx, conf.UnknownError)
		ctx.Abort()
		return
	}

	accessKey := &account.UserAccessKeyBean{
		UserId: userBean.ID,
		Key:    strings.ReplaceAll(uuid.New().String(), "-", ""),   // uuid随机产生
		Secret: hex.EncodeToString(algorithmutils.RandomHmacKey()), // 随机产生秘钥
	}

	err := models.CreateBaseRecord(accessKey)
	if err != nil {
		logger.Info("GetUserByEmail sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, accessKey)
}
