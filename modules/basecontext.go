package modules

import (
	"encoding/json"
	"net/http"
	"tpayment/conf"
	"tpayment/pkg/tlog"

	"github.com/labstack/echo"
)

func BaseError(context echo.Context, err conf.ResultCode) {
	logger := tlog.GetLogger(context)

	baseResp := &BaseResponse{
		ErrorCode:    err,
		ErrorMessage: err.String(),
	}

	resp, _ := json.Marshal(baseResp)
	logger.Info("response->", string(resp))

	// nolint
	_ = context.JSON(http.StatusBadRequest, baseResp)
}

func BaseSuccess(context echo.Context, data interface{}) {
	logger := tlog.GetLogger(context)

	baseResponse := BaseResponse{
		ErrorCode: conf.SUCCESS,
		Data:      data,
	}
	resp, _ := json.Marshal(baseResponse)
	logger.Info("response->", string(resp))

	// nolint
	_ = context.JSON(http.StatusOK, baseResponse)
}

type BaseResponse struct {
	ErrorCode    conf.ResultCode `json:"code"`
	ErrorMessage string          `json:"msg,omitempty"`
	Data         interface{}     `json:"data,omitempty"`
}

type BaseIDRequest struct {
	ID uint `json:"id"`
}

type BaseQueryRequest struct {
	MerchantId uint              `json:"merchant_id,omitempty"`
	AgencyId   uint              `json:"agency_id,omitempty"`
	DeviceId   uint              `json:"device_id,omitempty"`
	AppId      uint              `json:"app_id,omitempty"`
	BatchId    uint              `json:"batch_id,omitempty"`
	Offset     uint              `json:"offset"`
	Limit      uint              `json:"limit"`
	Filters    map[string]string `json:"filters,omitempty"`
}

type BaseQueryResponse struct {
	Total uint `json:"total"`

	Data interface{} `json:"data"`
}
