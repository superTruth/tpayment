package router

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/constant"
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
			log := new(tlog.Logger)
			log.Init(requestId)
			ctx.Set(constant.LOG, log)
			defer log.Destroy()

			return handlerFunc(ctx)
		}
	}
}

func AuthHandle() echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if ctx.Request().RequestURI == conf.UrlLogin {  // 唯一的登录功能不需要token
				return handlerFunc(ctx)
			}

			userBean, _, err := user.Auth(ctx, conf.HeaderTagToken)
			if err != nil {   // 数据库出错
				modules.BaseError(ctx, conf.DBError)
				return err
			}

			if userBean == nil {  // token验证失败
				modules.BaseError(ctx, conf.TokenInvalid)
				return err
			}

			ctx.Set(conf.ContextTagUser, userBean)

			// 判断需要store id部分，普通用户，除了用户操作，其他操作都必须带store id


			return handlerFunc(ctx)
		}
	}
}
