package user

import (
	"tpayment/conf"
	"tpayment/internal/basekey"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

// 登录
func LoginHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(LoginRequest)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 创建 或者 更新  token记录
	accountBean, err := account.GetUserByEmail(models.DB(), ctx, req.Email)
	if err != nil {
		logger.Error("GetUserByEmail sql error->", err.Error())

		modules.BaseError(ctx, conf.DBError)
		return
	}

	if accountBean == nil { // 用户不存在的情况
		logger.Info("user not fund")
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	// 验证账号是否激活
	if !accountBean.Active {
		logger.Info("user not active")
		modules.BaseError(ctx, conf.UserNotActive)
		return
	}

	// 密码进行hash
	req.Pwd = basekey.Hash([]byte(req.Pwd))

	// 验证密码是否正确
	if accountBean.Pwd != req.Pwd {
		logger.Info("pwd error")
		modules.BaseError(ctx, conf.ValidateError)
		return
	}

	// 验证App id
	appBean, err := account.GetAppIdByAppID(models.DB(), ctx, req.AppId)
	if err != nil {
		logger.Error("GetAppIdByAppID sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	if appBean == nil { // App id不存在
		logger.Info("app not fund")
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	if appBean.AppSecret != req.AppSecret {
		logger.Info("app secret error")
		modules.BaseError(ctx, conf.ValidateError)
		return
	}

	// 查看是否已经存在这个账号的token，如果已经存在，直接update，如果不存在需要create
	tokenBean, err := account.GetTokenByUserId(models.DB(), ctx, accountBean.ID, appBean.ID)
	if err != nil {
		logger.Error("GetTokenByUserId sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	if tokenBean == nil { // 如果不存在需要create
		tokenBean = &account.TokenBean{
			UserId: accountBean.ID,
			AppId:  appBean.ID,
			Token:  uuid.New().String(),
		}
		if err = models.CreateBaseRecord(tokenBean); err != nil {
			logger.Error("CreateBaseRecord sql error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return
		}
	} else {
		tokenBean.Token = uuid.New().String()
		if err = models.UpdateBaseRecord(tokenBean); err != nil {
			logger.Error("UpdateBaseRecord sql error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return
		}
	}

	// 拼接response数据
	ret := &LoginResponse{
		UserID: accountBean.ID,
		Token:  tokenBean.Token,
		Role:   accountBean.Role,
		Name:   accountBean.Name,
		Email:  accountBean.Email,
	}

	modules.BaseSuccess(ctx, ret)
}
