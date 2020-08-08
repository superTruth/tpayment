package router

import (
	"errors"
	"net/http/httputil"
	"strings"
	"tpayment/conf"
	"tpayment/constant"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/modules/user"
	"tpayment/pkg/tlog"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

func PreHandle() echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// 生成request ID
			requestId := uuid.New().String()
			ctx.Set(constant.REQUEST_ID, requestId)

			ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
			ctx.Response().Header().Set("content-type", "application/json")

			// 生成log
			logger := new(tlog.Logger)
			logger.Init(requestId)
			tlog.SetLogger(ctx, logger)
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
			if userBean.Role == string(conf.RoleUser) { // 机器人和系统管理员不需要验证
				_, agencyBean, err := agency.QueryAgencyRecord(models.DB(), ctx, 0, 1000, nil)
				if err != nil {
					logger.Error("QueryAgencyRecord db error->", err.Error())
					modules.BaseError(ctx, conf.DBError)
					return err
				}
				ctx.Set(conf.ContextTagAgency, agencyBean)
			} else {
				ctx.Set(conf.ContextTagAgency, make([]*agency.Agency, 0))
			}

			return handlerFunc(ctx)
		}
	}
}

func PermissionFilter() echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			logger := tlog.GetLogger(ctx)

			if ctx.Request().RequestURI == conf.UrlAccountLogin ||
				ctx.Request().RequestURI == conf.UrlAccountRegister ||
				strings.Contains(ctx.Request().RequestURI, "/payment/account/active") { // 唯一的登录功能不需要token
				return handlerFunc(ctx)
			}

			userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)

			// 管理员直行
			if userBean.Role == string(conf.RoleAdmin) {
				return handlerFunc(ctx)
			}

			// 需要管理员权限的
			if ctx.Request().RequestURI == conf.UrlAgencyAdd || // 所有的机构操作只能管理员
				ctx.Request().RequestURI == conf.UrlAgencyUpdate ||
				ctx.Request().RequestURI == conf.UrlAgencyDelete ||
				strings.Contains(ctx.Request().RequestURI, "/payment/agency_device") {
				logger.Warn("no agency permission")
				modules.BaseError(ctx, conf.NoPermission)
				return errors.New("not admin")
			}

			// 需要管理员或者机构管理员的权限的
			agencys := ctx.Get(conf.ContextTagAgency).([]*agency.Agency)

			if ctx.Request().RequestURI == conf.UrlAccountAdd ||
				ctx.Request().RequestURI == conf.UrlAccountDelete ||
				ctx.Request().RequestURI == conf.UrlAccountQuery ||
				ctx.Request().RequestURI == conf.UrlMerchantAdd ||
				ctx.Request().RequestURI == conf.UrlMerchantUpdate ||
				ctx.Request().RequestURI == conf.UrlMerchantQuery ||
				strings.Contains(ctx.Request().RequestURI, "/payment/tms") {

				if userBean.Role == string(conf.RoleUser) { // 管理员，不需要过滤机构
					if len(agencys) == 0 {
						logger.Info("not admin or not agency")
						modules.BaseError(ctx, conf.NoPermission)
						return errors.New("not admin or not agency")
					}
				}
			}

			return handlerFunc(ctx)
		}
	}
}
