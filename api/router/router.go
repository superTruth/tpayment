package router

import (
	"github.com/labstack/echo"
	"net/http"
	"tpayment/conf"
	"tpayment/modules/merchant"
	"tpayment/modules/merchant/associate"
	"tpayment/modules/user"
)

func Init() (*echo.Echo, error) {
	e := echo.New()

	e.Use(PreHandle())  // 前置过滤
	e.Use(AuthHandle()) // 授权过滤

	e.GET("/", func(context echo.Context) error {
		return context.String(http.StatusOK, "Hello, world!")
	})

	//UrlAccountAdd    = "/payment/account/add"
	//UrlAccountDelete = "/payment/account/delete"
	//UrlAccountUpdate = "/payment/account/update"
	//UrlAccountQuery  = "/payment/account/query"

	e.POST(conf.UrlAccountLogin, user.LoginHandle)   // 登录账号
	e.POST(conf.UrlAccountLogout, user.LogoutHandle) // 登出账号
	e.POST(conf.UrlAccountAdd, user.AddHandle)       // 新增账号
	e.POST(conf.UrlAccountDelete, user.DeleteHandle) // 删除账号
	e.POST(conf.UrlAccountUpdate, user.UpdateHandle) // 更新账号
	e.POST(conf.UrlAccountQuery, user.QueryHandle)   // 查找账号

	e.POST(conf.UrlMerchantAdd, merchant.AddHandle)       // 新增商户
	e.POST(conf.UrlMerchantUpdate, merchant.UpdateHandle) // 更新商户
	e.POST(conf.UrlMerchantQuery, merchant.QueryHandle)   // 查找商户

	e.POST(conf.UrlMerchantAssociateAdd, associate.AddHandle)       // 新增商户
	e.POST(conf.UrlMerchantAssociateDelete, associate.DeleteHandle) // 删除账号
	e.POST(conf.UrlMerchantAssociateQuery, associate.QueryHandle)   // 查找商户
	e.POST(conf.UrlQueryUserInMerchantQuery, associate.QueryUserInMerchantHandle)
	//e.POST("/batchupdate", BatchUpdate) // 批量升级功能
	//
	//e.POST("/uploadfilerequest", RequestUploadFileUrl) // 申请文件上传
	//
	//e.POST("/createfile", CreateFile) // 创建新文件

	return e, nil
}
