package errorcode

type CommonError struct {
	Msg  string `json:"msg,omitempty"`
	Code int    `json:"code,omitempty"`
}

// 根据error生成CommonError
func GenerateError(err error) *CommonError {

	commonErr := new(CommonError)
	commonErr.Code = -1
	commonErr.Msg = err.Error()

	return commonErr
}

// 根据error生成CommonError
func GenerateErrorFromStr(msg string) *CommonError {

	commonErr := new(CommonError)
	commonErr.Code = -1
	commonErr.Msg = msg

	return commonErr
}

//
var ERROR_PARAMETERS_ERROR = &CommonError{Code: 2, Msg: "parameters error"}
var ERROR_CANT_FIND_RECORD = &CommonError{Code: 3, Msg: "can't find record"}
