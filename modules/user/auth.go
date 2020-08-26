package user

import (
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context, token string) (*account.UserBean, *account.AppIdBean, error) {
	logger := tlog.GetLogger(ctx)
	// 创建 或者 更新  token记录
	tokenBean, err := account.GetTokenBeanByToken(models.DB(), ctx, token)
	if err != nil {
		logger.Warn("GetTokenBeanByToken fail->", err.Error())
		return nil, nil, err
	}

	if tokenBean == nil { // 没有对应的token记录
		return nil, nil, nil
	}

	//
	accountBean, err := account.GetUserById(models.DB(), ctx, tokenBean.UserId)
	if err != nil {
		logger.Warn("GetUserById fail->", err.Error())
		return nil, nil, err
	}

	appBean, err := account.GetAppIdByID(models.DB(), ctx, tokenBean.AppId)
	if err != nil {
		logger.Warn("GetUserById fail->", err.Error())
		return nil, nil, err
	}

	return accountBean, appBean, nil
}
