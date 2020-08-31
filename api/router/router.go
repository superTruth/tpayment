package router

import (
	"net/http"
	"tpayment/api/router/middleware"
	"tpayment/conf"
	"tpayment/modules/agency"
	"tpayment/modules/agency/acquirer"
	"tpayment/modules/agency/agencydevice"
	associate2 "tpayment/modules/agency/associate"
	"tpayment/modules/agency/payment"
	"tpayment/modules/fileupload"
	"tpayment/modules/merchant"
	"tpayment/modules/merchant/associate"
	"tpayment/modules/merchant/merchantdevice"
	"tpayment/modules/merchant/merchantdevicepayment"
	"tpayment/modules/tms/app"
	"tpayment/modules/tms/appfile"
	"tpayment/modules/tms/appinbatchupdate"
	"tpayment/modules/tms/appindevice"
	"tpayment/modules/tms/batchupdate"
	"tpayment/modules/tms/clientapi"
	"tpayment/modules/tms/device"
	"tpayment/modules/tms/devicemodel"
	"tpayment/modules/tms/devicetag"
	"tpayment/modules/tms/uploadfile"
	"tpayment/modules/user"
	"tpayment/modules/user/accesskey"

	"github.com/gin-gonic/gin"
)

func Init() (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()

	e.Use(
		middleware.Logger,
		middleware.NewRecovery,
		//middleware.Cors,
		middleware.NewCors(),
		middleware.AuthHandle,
		middleware.PermissionFilter,
	)

	e.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello, world!")
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
	e.POST(conf.UrlAccountAccessAdd, accesskey.AddHandle)
	e.POST(conf.UrlAccountAccessQuery, accesskey.QueryHandle)
	e.POST(conf.UrlAccountAccessDelete, accesskey.DeleteHandle)

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

	e.POST(conf.UrlAgencyDeviceAdd, agencydevice.AddHandle)              // 添加机构设备
	e.POST(conf.UrlAgencyDeviceDelete, agencydevice.DeleteHandle)        // 删除
	e.POST(conf.UrlAgencyDeviceQuery, agencydevice.QueryAssociateHandle) // 查询

	e.POST(conf.UrlAgencyAcquirerAdd, acquirer.AddHandle)       // 添加acquirer
	e.POST(conf.UrlAgencyAcquirerUpdate, acquirer.UpdateHandle) // 更新
	e.POST(conf.UrlAgencyAcquirerQuery, acquirer.QueryHandle)   // 查找
	e.POST(conf.UrlAgencyAcquirerDelete, acquirer.DeleteHandle) // 删除

	e.POST(conf.UrlMerchantAdd, merchant.AddHandle)                     // 新增商户
	e.POST(conf.UrlMerchantUpdate, merchant.UpdateHandle)               // 更新
	e.POST(conf.UrlMerchantQuery, merchant.QueryHandle)                 // 查找
	e.POST(conf.UrlMerchantDelete, merchant.DeleteHandle)               // 删除
	e.POST(conf.UrlMerchantInAgencyQuery, merchant.InAgencyQueryHandle) // 查询Agency下面有多少商户

	e.POST(conf.UrlMerchantAssociateAdd, associate.AddHandle)       // 新增商户员工
	e.POST(conf.UrlMerchantAssociateDelete, associate.DeleteHandle) // 删除商户员工
	e.POST(conf.UrlMerchantAssociateQuery, associate.QueryHandle)   // 查询商户员工
	e.POST(conf.UrlMerchantAssociateUpdate, associate.UpdateHandle) // 查询商户员工

	e.POST(conf.UrlMerchantDeviceAdd, merchantdevice.AddHandle)       // 新增商户设备
	e.POST(conf.UrlMerchantDeviceDelete, merchantdevice.DeleteHandle) // 删除关联
	e.POST(conf.UrlMerchantDeviceUpdate, merchantdevice.UpdateHandle) // 更新
	e.POST(conf.UrlMerchantDeviceQuery, merchantdevice.QueryHandle)   // 查询

	e.POST(conf.UrlMerchantDevicePaymentAdd, merchantdevicepayment.AddHandle)       // 新增商户设备支付参数
	e.POST(conf.UrlMerchantDevicePaymentDelete, merchantdevicepayment.DeleteHandle) // 删除
	e.POST(conf.UrlMerchantDevicePaymentUpdate, merchantdevicepayment.UpdateHandle) // 更新
	e.POST(conf.UrlMerchantDevicePaymentQuery, merchantdevicepayment.QueryHandle)   // 查询

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

	e.POST(conf.UrlTmsAppFileAdd, appfile.AddHandle)       // 添加app fileutils
	e.POST(conf.UrlTmsAppFileDelete, appfile.DeleteHandle) // 删除
	e.POST(conf.UrlTmsAppFileUpdate, appfile.UpdateHandle) // 更新
	e.POST(conf.UrlTmsAppFileQuery, appfile.QueryHandle)   // 查询

	e.POST(conf.UrlTmsTagAdd, devicetag.AddHandle)       // 添加device tag
	e.POST(conf.UrlTmsTagDelete, devicetag.DeleteHandle) // 删除
	e.POST(conf.UrlTmsTagUpdate, devicetag.UpdateHandle) // 更新
	e.POST(conf.UrlTmsTagQuery, devicetag.QueryHandle)   // 查询

	e.POST(conf.UrlTmsModelAdd, devicemodel.AddHandle)       // 添加device model
	e.POST(conf.UrlTmsModelDelete, devicemodel.DeleteHandle) // 删除
	e.POST(conf.UrlTmsModelUpdate, devicemodel.UpdateHandle) // 更新
	e.POST(conf.UrlTmsModelQuery, devicemodel.QueryHandle)   // 查询

	e.POST(conf.UrlTmsBatchUpdateAdd, batchupdate.AddHandle)           // 添加batch update
	e.POST(conf.UrlTmsBatchUpdateDelete, batchupdate.DeleteHandle)     // 删除
	e.POST(conf.UrlTmsBatchUpdateUpdate, batchupdate.UpdateHandle)     // 更新
	e.POST(conf.UrlTmsBatchUpdateQuery, batchupdate.QueryHandle)       // 查询
	e.POST(conf.UrlTmsBatchUpdateStartHandle, batchupdate.StartHandle) // 处理

	e.POST(conf.UrlTmsAppInBatchUpdateAdd, appinbatchupdate.AddHandle)       // 添加app in batch update
	e.POST(conf.UrlTmsAppInBatchUpdateDelete, appinbatchupdate.DeleteHandle) // 删除
	e.POST(conf.UrlTmsAppInBatchUpdateUpdate, appinbatchupdate.UpdateHandle) // 更新
	e.POST(conf.UrlTmsAppInBatchUpdateQuery, appinbatchupdate.QueryHandle)   // 查询

	e.POST(conf.UrlTmsUploadFileAdd, uploadfile.AddHandle)       // 添加app in batch update
	e.POST(conf.UrlTmsUploadFileDelete, uploadfile.DeleteHandle) // 删除
	e.POST(conf.UrlTmsUploadFileQuery, uploadfile.QueryHandle)   // 查询
	//e.POST(conf.UrlTmsUploadFileUpdate, uploadfile.UpdateHandle) // 更新

	e.POST(conf.UrlTmsHeartBeat, clientapi.HearBeat) // 客户端心跳逻辑

	e.POST(conf.UrlFileAdd, fileupload.RequestUploadFileUrl) // 创建文件

	return e, nil
}
