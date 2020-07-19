package conf

type ResultCode string

const (
	SUCCESS                ResultCode = "00"
	UnknownError           ResultCode = "99"
	ParameterError         ResultCode = "03"
	TokenInvalid           ResultCode = "04"
	DBError                ResultCode = "05"
	RecordNotFund          ResultCode = "06"
	ValidateError          ResultCode = "07"
	RecordAlreadyExist     ResultCode = "08"
	AuthFail               ResultCode = "09"
	NoPermission           ResultCode = "10"
	NeedTokenInHeader      ResultCode = "11"
	SendEmailFail          ResultCode = "12"
	UserNotActive          ResultCode = "13"
	UserCanOnlyInOneAgency ResultCode = "14"
	AdminCantAssociate     ResultCode = "15"
)

var ResultCodeText = map[ResultCode]string{
	SUCCESS:                "success",
	ParameterError:         "parameter error",
	TokenInvalid:           "token invalid",
	DBError:                "internal error 05",
	RecordNotFund:          "record not fund",
	RecordAlreadyExist:     "record already exist",
	AuthFail:               "auth fail",
	NoPermission:           "no permission",
	NeedTokenInHeader:      "need token in header",
	SendEmailFail:          "send email fail",
	UserNotActive:          "user not active",
	UserCanOnlyInOneAgency: "user can only in one agency",
	AdminCantAssociate:     "admin cant associate",
}

func (this ResultCode) String() string {
	return ResultCodeText[this]
}
