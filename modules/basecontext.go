package modules

import (
	"github.com/labstack/echo"
	"tpayment/conf"
)

func BaseError(context echo.Context, err conf.ResultCode) {

}

func BaseSuccess(context echo.Context, data interface{}) {

}

type BaseResponse struct {
	ErrorCode    conf.ResultCode `json:"error_code"`
	ErrorMessage string          `json:"error_message"`
}

type BaseIDRequest struct {
	ID uint `json:"id"`
}

type BaseQueryRequest struct {
	Offset  uint              `json:"offset"`
	Limit   uint              `json:"limit"`
	Filters map[string]string `json:"filters"`
}

type BaseQueryResponse struct {
	Total uint `json:"total"`

	Data []map[string]interface{} `json:"data"`
}
