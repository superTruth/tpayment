package conf

const HeaderTagRequestId = "X-Request-Id"
const HeaderTagToken = "X-Access-Token"
const HeaderTagAccessKey = "X-Access-Key"
const HeaderTagAccessHash = "X-Access-Hash"

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
	DbType           = "mysql"
	DbConnectInfoTag = "TPAYMENT_DB_CONFIG"
)

const (
	RebootModeNever      = "never"
	RebootModeEveryDay   = "every_day"
	RebootModeEveryWeek  = "every_week"
	RebootModeEveryMonth = "every_month"
)

const (
	TmsStatusPendingInstall     = "pending_install"
	TmsStatusInstalled          = "installed"
	TmsStatusPendingUninstalled = "pending_uninstall"
	TmsStatusWarningInstalled   = "warning_installed"
)

const (
	AppFileStatusPending  = "pending"
	AppFileStatusDecoding = "decoding"
	AppFileStatusDone     = "done"
	AppFileStatusFail     = "fail"
)
