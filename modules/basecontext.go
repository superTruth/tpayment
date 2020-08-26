package modules

import (
	"net/http"
	"tpayment/conf"

	"github.com/gin-gonic/gin"
)

func BaseError(context *gin.Context, err conf.ResultCode) {
	baseResp := &BaseResponse{
		ErrorCode:    err,
		ErrorMessage: err.String(),
	}

	context.JSON(http.StatusBadRequest, baseResp)
}

func BaseSuccess(context *gin.Context, data interface{}) {

	baseResponse := BaseResponse{
		ErrorCode: conf.SUCCESS,
		Data:      data,
	}

	context.JSON(http.StatusOK, baseResponse)
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
