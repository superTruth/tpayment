package conf

const HeaderTagRequestId = "RequestId"
const HeaderTagToken = "Token"

const ContextTagLog = "Log"
const ContextTagUser = "User"

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleUser    UserRole = "user"
	RoleMachine UserRole = "machine"
)

const MaxQueryCount = 100

const (
	UrlAccountLogin  = "/payment/account/login"
	UrlAccountLogout = "/payment/account/logout"
	UrlAccountAdd    = "/payment/account/add"
	UrlAccountDelete = "/payment/account/delete"
	UrlAccountUpdate = "/payment/account/update"
	UrlAccountQuery  = "/payment/account/query"

	UrlAgencyAdd    = "/payment/agency/add"
	UrlAgencyUpdate = "/payment/agency/update"
	UrlAgencyQuery  = "/payment/agency/query"

	UrlMerchantAdd    = "/payment/merchant/add"
	UrlMerchantUpdate = "/payment/merchant/update"
	UrlMerchantQuery  = "/payment/merchant/query"

	UrlMerchantAssociateAdd     = "/payment/merchant_associate/add"
	UrlMerchantAssociateDelete  = "/payment/merchant_associate/delete"
	UrlMerchantAssociateQuery   = "/payment/merchant_associate/query"
	UrlQueryUserInMerchantQuery = "/payment/merchant_associate/queryuserinmerchant"
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
