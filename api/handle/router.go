package router

import (
	"github.com/labstack/echo"
	"net/http"
	"tpayment/conf"
	"tpayment/modules/user"
)

func Init() (*echo.Echo, error) {
	e := echo.New()

	e.Use(PreHandle())  // 前置过滤
	e.Use(AuthHandle()) // 授权过滤

	e.GET("/", func(context echo.Context) error {
		return context.String(http.StatusOK, "Hello, world!")
	})

	e.POST(conf.UrlLogin, user.LoginHandle)   // 登录
	e.POST(conf.UrlLogout, user.LogoutHandle) // 登出

	//e.POST("/batchupdate", BatchUpdate) // 批量升级功能
	//
	//e.POST("/uploadfilerequest", RequestUploadFileUrl) // 申请文件上传
	//
	//e.POST("/createfile", CreateFile) // 创建新文件

	return e, nil
}
