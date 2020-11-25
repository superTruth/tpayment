package conf

type ResultCode string

const (
	Success                 ResultCode = "00"
	UnknownError            ResultCode = "99"
	PanicError              ResultCode = "98"
	ParameterError          ResultCode = "03"
	TokenInvalid            ResultCode = "04"
	DBError                 ResultCode = "05"
	RecordNotFund           ResultCode = "06"
	ValidateError           ResultCode = "07"
	RecordAlreadyExist      ResultCode = "08"
	AuthFail                ResultCode = "09"
	NoPermission            ResultCode = "10"
	NeedTokenInHeader       ResultCode = "11"
	SendEmailFail           ResultCode = "12"
	UserNotActive           ResultCode = "13"
	UserCanOnlyInOneAgency  ResultCode = "14"
	AdminCantAssociate      ResultCode = "15"
	DataIsUsing             ResultCode = "16"
	NotSupport              ResultCode = "17"
	DecodeError             ResultCode = "18"
	DecodeCardBrandError    ResultCode = "19"
	DecodeQRError           ResultCode = "20"
	NoPaymentProcessRule    ResultCode = "21"
	ProcessRuleSettingError ResultCode = "22"
	TIDIsBusy               ResultCode = "23"
	NoAvailableTID          ResultCode = "24"
	Reversal                ResultCode = "25"
	CantReachAcquirer       ResultCode = "26"
	RejectByAcquirer        ResultCode = "27"
)

var ResultCodeText = map[ResultCode]string{
	Success:                 "success",
	ParameterError:          "parameter error",
	UnknownError:            "internal error 99",
	PanicError:              "internal error 98",
	TokenInvalid:            "token invalid",
	ValidateError:           "validate error",
	DBError:                 "internal error 05",
	RecordNotFund:           "record not fund",
	RecordAlreadyExist:      "record already exist",
	AuthFail:                "auth fail",
	NoPermission:            "no permission",
	NeedTokenInHeader:       "need token in header",
	SendEmailFail:           "send email fail",
	UserNotActive:           "user not active",
	UserCanOnlyInOneAgency:  "user can only in one agency",
	AdminCantAssociate:      "admin cant associate",
	DataIsUsing:             "data in using",
	NotSupport:              "not support",
	DecodeError:             "decode error",
	DecodeCardBrandError:    "decode card brand error",
	DecodeQRError:           "decode qr code error",
	NoPaymentProcessRule:    "no payment process rule",
	ProcessRuleSettingError: "process rule setting error",
	TIDIsBusy:               "TID is busy",
	NoAvailableTID:          "no available TID",
	Reversal:                "need reversal",
	CantReachAcquirer:       "can't reach acquirer",
	RejectByAcquirer:        "reject by acquirer",
}

func (this ResultCode) String() string {
	return ResultCodeText[this]
}
