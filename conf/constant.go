package conf

const HeaderTagRequestId = "X-Request-Id"
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
	DbType           = "mysql"
	DbConnectInfoTag = "TPAYMENT_DB_CONFIG"
)

const (
	RebootModeNever      = "Never"
	RebootModeEveryDay   = "Every Day"
	RebootModeEveryWeek  = "Every Week"
	RebootModeEveryMonth = "Every Month"
)

const (
	TmsStatusPendingInstall     = "pending install"
	TmsStatusInstalled          = "installed"
	TmsStatusPendingUninstalled = "pending uninstall"
	TmsStatusWarningInstalled   = "warning installed"
)

const (
	AppFileStatusPending  = "pending"
	AppFileStatusDecoding = "decoding"
	AppFileStatusDone     = "done"
	AppFileStatusFail     = "fail"
)
