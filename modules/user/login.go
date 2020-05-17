package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"
)

// 登录
func LoginHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(LoginRequest)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 创建 或者 更新  token记录
	accountBean, err := account.GetUserByEmail(req.Email)
	if err != nil {
		logger.Error("GetUserByEmail sql error->", err.Error())

		modules.BaseError(ctx, conf.DBError)
		return err
	}

	if accountBean == nil { // 用户不存在的情况
		logger.Info("user not fund")
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}

	// 验证密码是否正确
	if accountBean.Pwd != req.Pwd {
		logger.Info("pwd error")
		modules.BaseError(ctx, conf.ValidateError)
		return errors.New("pwd error")
	}

	// 验证App id
	appBean, err := account.GetAppIdByAppID(req.AppId)
	if err != nil {
		logger.Error("GetAppIdByAppID sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	if appBean == nil {  // App id不存在
		logger.Info("app not fund")
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}

	if appBean.AppSecret != req.AppSecret {
		logger.Info("app secret error")
		modules.BaseError(ctx, conf.ValidateError)
		return errors.New("app secret error")
	}

	// 查看是否已经存在这个账号的token，如果已经存在，直接update，如果不存在需要create
	tokenBean, err := account.GetTokenByUserId(accountBean.ID, appBean.ID)
	if err != nil {
		logger.Error("GetTokenByUserId sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	if tokenBean == nil {  // 如果不存在需要create
		tokenBean = &account.TokenBean{
			UserId: accountBean.ID,
			AppId:  appBean.ID,
			Token:  uuid.New().String(),
		}
		if err = models.CreateBaseRecord(tokenBean); err != nil {
			logger.Error("CreateBaseRecord sql error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return err
		}
	} else {
		tokenBean.Token = uuid.New().String()
		if err = models.UpdateBaseRecord(tokenBean); err != nil {
			logger.Error("UpdateBaseRecord sql error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return err
		}
	}

	// 拼接response数据
	ret := &LoginResponse{
		BaseResponse: modules.BaseResponse{
			ErrorCode: conf.SUCCESS,
		},
		Token: tokenBean.Token,
		Role:  accountBean.Role,
		Name:  accountBean.Name,
		Email: accountBean.Email,
	}

	modules.BaseSuccess(ctx, ret)

	return nil
}
