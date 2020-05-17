package conf

const HeaderTagRequestId = "RequestId"
const HeaderTagToken = "Token"
const HeaderTagStoreId = "StoreID"

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
	UrlLogin  = "/payment/account/login"
	UrlLogout = "/payment/account/logout"
)

const (
	DbType = "mysql"
	DbConnectInfoTag = "TPAYMENT_DB_CONFIG"
)