package user

import (
	"encoding/base64"
	"errors"
	"github.com/go-gomail/gomail"
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"
)

func RegisterHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(AddUserRequest)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 查询是否已经存在的账号
	user, err := account.GetUserByEmail(req.Email)
	if err != nil {
		logger.Info("GetUserByEmail sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if user == nil {
		// 生成新账号
		user = &account.UserBean{
			Email: req.Email,
			Pwd:   req.Pwd,
			Name:  req.Name,
		}
		err = models.CreateBaseRecord(user)
		if err != nil {
			logger.Info("CreateBaseRecord sql error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return err
		}

	} else {
		if user.Active {
			logger.Warn(conf.RecordAlreadyExist.String())
			modules.BaseError(ctx, conf.RecordAlreadyExist)
			return errors.New(conf.RecordAlreadyExist.String())
		}
	}

	// 发送email验证
	err = sendActiveEmail(user.Email)
	if err != nil {
		logger.Warn("sendActiveEmail error->", err.Error())
		modules.BaseError(ctx, conf.SendEmailFail)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}

// 发送激活邮件
func sendActiveEmail(email string) error {
	m := gomail.NewMessage()

	m.SetHeader("To", email)
	m.SetAddressHeader("From", conf.GetConfigData().EmailUserAccount, conf.GetConfigData().EmailUserName)
	m.SetHeader("Subject", "Active Email")

	body := "Active here <a href = " + conf.GetConfigData().Domain + "payment/account/active/" + base64.StdEncoding.EncodeToString([]byte(email)) + ">Click</a><br>"
	m.SetBody("text/html", body)

	d := gomail.NewDialer(conf.GetConfigData().EmailHost, conf.GetConfigData().EmailHostPort,
		conf.GetConfigData().EmailUserAccount, conf.GetConfigData().EmailUserPwd)
	err := d.DialAndSend(m)

	if err != nil {
		return err
	}

	return nil
}

func ActiveHandel(ctx echo.Context) error {
	userId := ctx.Param("user")

	logger := tlog.GetLogger(ctx)

	email, err := base64.StdEncoding.DecodeString(userId)
	if err != nil {
		logger.Info("DecodeString email error->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 查询是否已经存在的账号
	user, err := account.GetUserByEmail(string(email))
	if err != nil {
		logger.Info("GetUserByEmail sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if user == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}

	// 激活账号
	user.Active = true
	err = models.UpdateBaseRecord(user)
	if err != nil {
		logger.Info("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
