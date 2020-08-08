package conf

const (
	UrlAccountLogin    = "/payment/account/login"
	UrlAccountLogout   = "/payment/account/logout"
	UrlAccountValidate = "/payment/account/validate"
	UrlAccountAdd      = "/payment/account/add"
	UrlAccountDelete   = "/payment/account/delete"
	UrlAccountUpdate   = "/payment/account/update"
	UrlAccountQuery    = "/payment/account/query"
	UrlAccountRegister = "/payment/account/register"
	UrlAccountActive   = "/payment/account/active/:user"

	UrlAgencyAdd    = "/payment/agency/add"
	UrlAgencyUpdate = "/payment/agency/update"
	UrlAgencyQuery  = "/payment/agency/query"
	UrlAgencyDelete = "/payment/agency/delete"

	UrlAgencyPaymentMethods = "/payment/agency_acquirer/payment_methods"
	UrlAgencyPaymentTypes   = "/payment/agency_acquirer/payment_types"
	UrlAgencyEntryTypes     = "/payment/agency_acquirer/entry_types"

	UrlAgencyAssociateAdd    = "/payment/agency_associate/add"
	UrlAgencyAssociateDelete = "/payment/agency_associate/delete"
	UrlAgencyAssociateQuery  = "/payment/agency_associate/query"

	UrlAgencyDeviceAdd    = "/payment/agency_device/add"
	UrlAgencyDeviceDelete = "/payment/agency_device/delete"
	UrlAgencyDeviceQuery  = "/payment/agency_device/query"

	UrlAgencyAcquirerAdd    = "/payment/agency_acquirer/add"
	UrlAgencyAcquirerUpdate = "/payment/agency_acquirer/update"
	UrlAgencyAcquirerQuery  = "/payment/agency_acquirer/query"
	UrlAgencyAcquirerDelete = "/payment/agency_acquirer/delete"

	UrlMerchantAdd           = "/payment/merchant/add"
	UrlMerchantUpdate        = "/payment/merchant/update"
	UrlMerchantQuery         = "/payment/merchant/query"
	UrlMerchantDelete        = "/payment/merchant/delete"
	UrlMerchantInAgencyQuery = "/payment/merchant_in_agency/query"

	UrlMerchantAssociateAdd    = "/payment/merchant_associate/add"
	UrlMerchantAssociateDelete = "/payment/merchant_associate/delete"
	UrlMerchantAssociateQuery  = "/payment/merchant_associate/query"
	UrlMerchantAssociateUpdate = "/payment/merchant_associate/update"

	UrlMerchantDeviceAdd    = "/payment/merchant_device/add"
	UrlMerchantDeviceDelete = "/payment/merchant_device/delete"
	UrlMerchantDeviceQuery  = "/payment/merchant_device/query"
	UrlMerchantDeviceUpdate = "/payment/merchant_device/update"

	UrlMerchantDevicePaymentAdd    = "/payment/merchant_device_payment/add"
	UrlMerchantDevicePaymentDelete = "/payment/merchant_device_payment/delete"
	UrlMerchantDevicePaymentUpdate = "/payment/merchant_device_payment/update"
	UrlMerchantDevicePaymentQuery  = "/payment/merchant_device_payment/query"

	UrlTmsDeviceDelete = "/payment/tms/agencydevice/delete"
	UrlTmsDeviceUpdate = "/payment/tms/agencydevice/update"
	UrlTmsDeviceQuery  = "/payment/tms/agencydevice/query"

	UrlTmsAppInDeviceAdd    = "/payment/tms/deviceapp/add"
	UrlTmsAppInDeviceDelete = "/payment/tms/deviceapp/delete"
	UrlTmsAppInDeviceUpdate = "/payment/tms/deviceapp/update"
	UrlTmsAppInDeviceQuery  = "/payment/tms/deviceapp/query"

	UrlTmsAppAdd    = "/payment/tms/app/add"
	UrlTmsAppDelete = "/payment/tms/app/delete"
	UrlTmsAppUpdate = "/payment/tms/app/update"
	UrlTmsAppQuery  = "/payment/tms/app/query"

	UrlTmsAppFileAdd    = "/payment/tms/app_file/add"
	UrlTmsAppFileDelete = "/payment/tms/app_file/delete"
	UrlTmsAppFileUpdate = "/payment/tms/app_file/update"
	UrlTmsAppFileQuery  = "/payment/tms/app_file/query"

	UrlTmsTagAdd    = "/payment/tms/tag/add"
	UrlTmsTagDelete = "/payment/tms/tag/delete"
	UrlTmsTagUpdate = "/payment/tms/tag/update"
	UrlTmsTagQuery  = "/payment/tms/tag/query"

	UrlTmsBatchUpdateAdd         = "/payment/tms/batchupdate/add"
	UrlTmsBatchUpdateDelete      = "/payment/tms/batchupdate/delete"
	UrlTmsBatchUpdateUpdate      = "/payment/tms/batchupdate/update"
	UrlTmsBatchUpdateQuery       = "/payment/tms/batchupdate/query"
	UrlTmsBatchUpdateStartHandle = "/payment/tms/batchupdate/starthandle"

	UrlTmsAppInBatchUpdateAdd    = "/payment/tms/appinbatchupdate/add"
	UrlTmsAppInBatchUpdateDelete = "/payment/tms/appinbatchupdate/delete"
	UrlTmsAppInBatchUpdateUpdate = "/payment/tms/appinbatchupdate/update"
	UrlTmsAppInBatchUpdateQuery  = "/payment/tms/appinbatchupdate/query"

	UrlTmsUploadFileAdd    = "/payment/tms/uploadfile/add"
	UrlTmsUploadFileDelete = "/payment/tms/uploadfile/delete"
	UrlTmsUploadFileUpdate = "/payment/tms/uploadfile/update"
	UrlTmsUploadFileQuery  = "/payment/tms/uploadfile/query"

	UrlTmsHeartBeat = "/payment/tms/heartbeat"

	UrlFileAdd = "/payment/file/add"
)