package user

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"tpayment/conf"
	"tpayment/internal/basekey"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
)

func RegisterHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(AddUserRequest)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 密码进行hash
	req.Pwd = basekey.Hash([]byte(req.Pwd))

	// 查询是否已经存在的账号
	user, err := account.GetUserByEmail(models.DB(), ctx, req.Email)
	if err != nil {
		logger.Info("GetUserByEmail sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if user == nil {
		// 如果不是机器的话，就直接是普通用户
		if req.Role != string(conf.RoleMachine) {
			req.Role = string(conf.RoleUser)
		}

		// 生成新账号
		user = &account.UserBean{
			Email:  req.Email,
			Pwd:    req.Pwd,
			Name:   req.Name,
			Role:   req.Role,
			Active: req.Role == string(conf.RoleMachine),
		}
		err = models.CreateBaseRecord(user)
		if err != nil {
			logger.Info("CreateBaseRecord sql error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return
		}

		// 机器人注册，直接成功
		if req.Role == string(conf.RoleMachine) {
			modules.BaseSuccess(ctx, nil)
			return
		}

	} else {
		if user.Role == string(conf.RoleMachine) { // 如果是机器人，查看是否登录过，如果未登录过，允许再次注册
			flag, err := account.TokenBeanDao.IsUserLogined(user.ID)
			if err != nil {
				logger.Info("IsUserLogined error->", err.Error())
				modules.BaseError(ctx, conf.DBError)
				return
			}
			if flag { // 如果已经登录过，则直接不允许重复注册
				logger.Warn(conf.RecordAlreadyExist.String())
				modules.BaseError(ctx, conf.RecordAlreadyExist)
				return
			}

			//
			logger.Warn("user multiple register->", user.ID)
			user.Pwd = req.Pwd
			if err = user.UpdatePwd(); err != nil {
				logger.Error("UpdatePwd fail->", err.Error())
				modules.BaseError(ctx, conf.RecordAlreadyExist)
				return
			}

			modules.BaseSuccess(ctx, nil)
			return
		}

		if user.Active {
			logger.Warn(conf.RecordAlreadyExist.String())
			modules.BaseError(ctx, conf.RecordAlreadyExist)
			return
		}
	}

	// 发送email验证
	err = sendActiveEmail(user.Email)
	if err != nil {
		logger.Warn("sendActiveEmail error->", err.Error())
		modules.BaseError(ctx, conf.SendEmailFail)
		return
	}

	modules.BaseSuccess(ctx, nil)
}

// 发送激活邮件
func sendActiveEmail(email string) error {
	m := gomail.NewMessage()

	m.SetHeader("To", email)
	m.SetAddressHeader("From", conf.GetConfigData().EmailUserAccount, conf.GetConfigData().EmailUserName)
	m.SetHeader("Subject", "Active Email")

	body := "Active here <a href = " + conf.GetConfigData().Domain + "payment/account/active/" + base64.StdEncoding.EncodeToString([]byte(email)) + ">Click</a><br>"
	m.SetBody("text/html", body)

	fmt.Println("Truth url->", body)

	d := gomail.NewDialer(conf.GetConfigData().EmailHost, conf.GetConfigData().EmailHostPort,
		conf.GetConfigData().EmailUserAccount, conf.GetConfigData().EmailUserPwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)

	if err != nil {
		return err
	}

	return nil
}

func ActiveHandel(ctx *gin.Context) {
	userId := ctx.Param("user")

	logger := tlog.GetLogger(ctx)

	email, err := base64.StdEncoding.DecodeString(userId)
	if err != nil {
		logger.Info("DecodeString email error->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 查询是否已经存在的账号
	user, err := account.GetUserByEmail(models.DB(), ctx, string(email))
	if err != nil {
		logger.Info("GetUserByEmail sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if user == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	// 激活账号
	user.Active = true
	err = models.UpdateBaseRecord(user)
	if err != nil {
		logger.Info("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	ctx.JSON(http.StatusOK, "Your Account:"+user.Email+" Active Success")
}
