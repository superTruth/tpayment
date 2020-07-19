package router

import (
	"github.com/labstack/echo"
	"net/http"
	"tpayment/conf"
	"tpayment/modules/agency"
	associate2 "tpayment/modules/agency/associate"
	"tpayment/modules/agency/payment"
	"tpayment/modules/merchant"
	"tpayment/modules/merchant/associate"
	"tpayment/modules/merchant/merchantdevice"
	"tpayment/modules/tms/app"
	"tpayment/modules/tms/appindevice"
	"tpayment/modules/tms/device"
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
	e.POST(conf.UrlAgencyDelete, agency.DeleteHandle) // 删除机构

	e.POST(conf.UrlAgencyPaymentMethods, payment.QueryPaymentMethodsHandle) // 查询支付方式
	e.POST(conf.UrlAgencyPaymentTypes, payment.QueryPaymentTypesHandle)     // 查询用卡方式
	e.POST(conf.UrlAgencyEntryTypes, payment.QueryEntryTypesHandle)         // 查询支付类型

	e.POST(conf.UrlAgencyAssociateAdd, associate2.AddHandle)              // 添加机构账户关联
	e.POST(conf.UrlAgencyAssociateDelete, associate2.DeleteHandle)        // 删除机构账户关联
	e.POST(conf.UrlAgencyAssociateQuery, associate2.QueryAssociateHandle) // 删除机构账户关联

	e.POST(conf.UrlMerchantAdd, merchant.AddHandle)       // 新增商户
	e.POST(conf.UrlMerchantUpdate, merchant.UpdateHandle) // 更新商户
	e.POST(conf.UrlMerchantQuery, merchant.QueryHandle)   // 查找商户

	e.POST(conf.UrlMerchantAssociateAdd, associate.AddHandle)       // 新增商户员工
	e.POST(conf.UrlMerchantAssociateDelete, associate.DeleteHandle) // 删除商户员工
	e.POST(conf.UrlMerchantAssociateQuery, associate.QueryHandle)   // 查询商户员工
	e.POST(conf.UrlMerchantAssociateUpdate, associate.UpdateHandle) // 查询商户员工

	e.POST(conf.UrlMerchantDevicePaymentAdd, merchantdevice.AddHandle)       // 新增商户设备
	e.POST(conf.UrlMerchantDevicePaymentDelete, merchantdevice.DeleteHandle) // 删除关联
	e.POST(conf.UrlMerchantDevicePaymentUpdate, merchantdevice.UpdateHandle) // 更新
	e.POST(conf.UrlMerchantDevicePaymentQuery, merchantdevice.QueryHandle)   // 查询

	e.POST(conf.UrlTmsDeviceDelete, device.DeleteHandle) // 删除关联
	e.POST(conf.UrlTmsDeviceUpdate, device.UpdateHandle) // 更新
	e.POST(conf.UrlTmsDeviceQuery, device.QueryHandle)   // 查询

	e.POST(conf.UrlTmsAppInDeviceAdd, appindevice.AddHandle)       // 添加app
	e.POST(conf.UrlTmsAppInDeviceDelete, appindevice.DeleteHandle) // 删除
	e.POST(conf.UrlTmsAppInDeviceUpdate, appindevice.UpdateHandle) // 更新
	e.POST(conf.UrlTmsAppInDeviceQuery, appindevice.QueryHandle)   // 查询

	e.POST(conf.UrlTmsAppAdd, app.AddHandle)       // 添加app
	e.POST(conf.UrlTmsAppDelete, app.DeleteHandle) // 删除
	e.POST(conf.UrlTmsAppUpdate, app.UpdateHandle) // 更新
	e.POST(conf.UrlTmsAppQuery, app.QueryHandle)   // 查询

	return e, nil
}
