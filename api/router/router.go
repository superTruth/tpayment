package router

import (
	"github.com/labstack/echo"
	"net/http"
	"tpayment/conf"
	"tpayment/modules/agency"
	associate2 "tpayment/modules/agency/associate"
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

	e.POST(conf.UrlAccountLogin, user.LoginHandle)       // 登录账号
	e.POST(conf.UrlAccountLogout, user.LogoutHandle)     // 登出账号
	e.POST(conf.UrlAccountValidate, user.ValidateHandle) // 验证账号
	e.POST(conf.UrlAccountAdd, user.AddHandle)           // 新增账号
	e.POST(conf.UrlAccountDelete, user.DeleteHandle)     // 删除账号
	e.POST(conf.UrlAccountUpdate, user.UpdateHandle)     // 更新账号
	e.POST(conf.UrlAccountQuery, user.QueryHandle)       // 查找账号
	e.POST(conf.UrlAccountRegister, user.RegisterHandle) // 注册账号
	e.GET(conf.UrlAccountActive, user.ActiveHandel)      // 激活账号

	e.POST(conf.UrlAgencyAdd, agency.AddHandle)       // 添加机构
	e.POST(conf.UrlAgencyUpdate, agency.UpdateHandle) // 更新机构
	e.POST(conf.UrlAgencyQuery, agency.QueryHandle)   // 查找机构

	e.POST(conf.UrlAgencyAssociateAdd, associate2.AddHandle) //

	e.POST(conf.UrlMerchantAdd, merchant.AddHandle)       // 新增商户
	e.POST(conf.UrlMerchantUpdate, merchant.UpdateHandle) // 更新商户
	e.POST(conf.UrlMerchantQuery, merchant.QueryHandle)   // 查找商户

	e.POST(conf.UrlMerchantAssociateAdd, associate.AddHandle)                     // 新增商户
	e.POST(conf.UrlMerchantAssociateDelete, associate.DeleteHandle)               // 删除关联
	e.POST(conf.UrlQueryUserInMerchantQuery, associate.QueryUserInMerchantHandle) // 查询

	return e, nil
}
