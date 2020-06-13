package router

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"net/http/httputil"
	"strings"
	"tpayment/conf"
	"tpayment/constant"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/modules/user"
	"tpayment/pkg/tlog"
)

func PreHandle() echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// 生成request ID
			requestId := uuid.New().String()
			ctx.Set(constant.REQUEST_ID, requestId)

			// 生成log
			logger := new(tlog.Logger)
			logger.Init(requestId)
			ctx.Set(constant.LOG, logger)
			defer logger.Destroy()

			content, _ := httputil.DumpRequest(ctx.Request(), true)
			logger.Info("request->", string(content))
			return handlerFunc(ctx)
		}
	}
}

func AuthHandle() echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			logger := tlog.GetLogger(ctx)

			if ctx.Request().RequestURI == conf.UrlAccountLogin ||
				ctx.Request().RequestURI == conf.UrlAccountRegister ||
				strings.Contains(ctx.Request().RequestURI, "/payment/account/active") { // 唯一的登录功能不需要token
				return handlerFunc(ctx)
			}

			// 判断token
			tokens := ctx.Request().Header[conf.HeaderTagToken]
			if len(tokens) == 0 {
				logger.Info("authHandle error->", conf.NeedTokenInHeader.String())
				modules.BaseError(ctx, conf.NeedTokenInHeader)
				return errors.New(conf.NeedTokenInHeader.String())
			}

			// 验证token的有效性
			userBean, _, err := user.Auth(ctx, tokens[0])
			if err != nil { // 数据库出错
				logger.Error("Auth db error->", err.Error())
				modules.BaseError(ctx, conf.DBError)
				return err
			}
			if userBean == nil { // token验证失败
				logger.Error("authHandle token not exist")
				modules.BaseError(ctx, conf.TokenInvalid)
				return err
			}

			ctx.Set(conf.ContextTagUser, userBean)

			// 查看是否是机构管理员
			if userBean.Role != string(conf.RoleUser) { // 机器人和系统管理员不需要验证
				_, agencyBean, err := agency.QueryAgencyRecord(models.DB(), ctx, 0, 1000, nil)
				if err != nil {
					logger.Error("QueryAgencyRecord db error->", err.Error())
					modules.BaseError(ctx, conf.DBError)
					return err
				}
				ctx.Set(conf.ContextTagAgency, agencyBean)
			}

			// 判断需要merchant id部分，普通用户，除了用户操作，其他操作都必须带merchant id

			return handlerFunc(ctx)
		}
	}
}
