package modules

import (
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"tpayment/conf"
	"tpayment/pkg/tlog"
)

func BaseError(context echo.Context, err conf.ResultCode) {
	logger := tlog.GetLogger(context)

	baseResp := &BaseResponse{
		ErrorCode:    err,
		ErrorMessage: err.String(),
	}

	resp, _ := json.Marshal(baseResp)
	logger.Info("response->", string(resp))

	context.JSON(http.StatusBadRequest, baseResp)
}

func BaseSuccess(context echo.Context, data interface{}) {
	logger := tlog.GetLogger(context)
	resp, _ := json.Marshal(data)
	logger.Info("response->", string(resp))

	context.JSON(http.StatusOK, data)
}

type BaseResponse struct {
	ErrorCode    conf.ResultCode `json:"error_code"`
	ErrorMessage string          `json:"error_message"`
}

type BaseIDRequest struct {
	ID uint `json:"id"`
}

type BaseQueryRequest struct {
	MerchantId uint              `json:"merchant_id"`
	Offset     uint              `json:"offset"`
	Limit      uint              `json:"limit"`
	Filters    map[string]string `json:"filters"`
}

type BaseQueryResponse struct {
	Total uint `json:"total"`

	Data interface{} `json:"data"`
}
