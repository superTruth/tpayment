package conf

const HeaderTagRequestId = "X-ACCESS-REQUEST-ID"
const HeaderTagToken = "Token"

const ContextTagLog = "Log"
const ContextTagUser = "User"
const ContextTagAgency = "Agency"

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleUser    UserRole = "user"
	RoleMachine UserRole = "machine"
)

const MaxQueryCount = 100

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

	UrlMerchantAdd    = "/payment/merchant/add"
	UrlMerchantUpdate = "/payment/merchant/update"
	UrlMerchantQuery  = "/payment/merchant/query"

	UrlMerchantAssociateAdd    = "/payment/merchant_associate/add"
	UrlMerchantAssociateDelete = "/payment/merchant_associate/delete"
	UrlMerchantAssociateQuery  = "/payment/merchant_associate/query"
	UrlMerchantAssociateUpdate = "/payment/merchant_associate/update"

	UrlMerchantDevicePaymentAdd    = "/payment/merchant_device_payment/add"
	UrlMerchantDevicePaymentDelete = "/payment/merchant_device_payment/delete"
	UrlMerchantDevicePaymentUpdate = "/payment/merchant_device_payment/update"
	UrlMerchantDevicePaymentQuery  = "/payment/merchant_device_payment/query"

	UrlTmsDeviceDelete = "/payment/tms/device/delete"
	UrlTmsDeviceUpdate = "/payment/tms/device/update"
	UrlTmsDeviceQuery  = "/payment/tms/device/query"

	UrlTmsAppInDeviceAdd    = "/payment/tms/deviceapp/add"
	UrlTmsAppInDeviceDelete = "/payment/tms/deviceapp/delete"
	UrlTmsAppInDeviceUpdate = "/payment/tms/deviceapp/update"
	UrlTmsAppInDeviceQuery  = "/payment/tms/deviceapp/query"

	UrlTmsAppAdd    = "/payment/tms/app/add"
	UrlTmsAppDelete = "/payment/tms/app/delete"
	UrlTmsAppUpdate = "/payment/tms/app/update"
	UrlTmsAppQuery  = "/payment/tms/app/query"
)

const (
	DbType           = "mysql"
	DbConnectInfoTag = "TPAYMENT_DB_CONFIG"
)

const (
	RebootModeNever    = "Never"
	RebootModeEveryDay = "Every Day"
)

const (
	TmsStatusPendingInstall     = "pending install"
	TmsStatusInstalled          = "installed"
	TmsStatusPendingUninstalled = "pending uninstall"
	TmsStatusWarningInstalled   = "warning installed"
)
