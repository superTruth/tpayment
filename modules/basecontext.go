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

	baseResponse := BaseResponse{
		ErrorCode: conf.SUCCESS,
		Data:      data,
	}
	resp, _ := json.Marshal(baseResponse)
	logger.Info("response->", string(resp))

	context.JSON(http.StatusOK, baseResponse)
}

type BaseResponse struct {
	ErrorCode    conf.ResultCode `json:"code"`
	ErrorMessage string          `json:"msg"`
	Data         interface{}     `json:"data"`
}

type BaseIDRequest struct {
	ID uint `json:"id"`
}

type BaseQueryRequest struct {
	MerchantId uint              `json:"merchant_id"`
	AgencyId   uint              `json:"agency_id"`
	Offset     uint              `json:"offset"`
	Limit      uint              `json:"limit"`
	Filters    map[string]string `json:"filters"`
}

type BaseQueryResponse struct {
	Total uint `json:"total"`

	Data interface{} `json:"data"`
}
