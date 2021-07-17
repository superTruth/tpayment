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
		ErrorCode: conf.Success,
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
	ID uint64 `json:"id"`
}

type BaseQueryRequest struct {
	MerchantId uint64            `json:"merchant_id,omitempty"`
	AgencyId   uint64            `json:"agency_id,omitempty"`
	DeviceId   uint64            `json:"device_id,omitempty"`
	DeviceSN   string            `json:"device_sn,omitempty"`
	AppId      uint64            `json:"app_id,omitempty"`
	BatchId    uint64            `json:"batch_id,omitempty"`
	TagId      uint64            `json:"tag_id,omitempty"`
	Offset     uint64            `json:"offset"`
	Limit      uint64            `json:"limit"`
	Filters    map[string]string `json:"filters,omitempty"`
}

type BaseQueryResponse struct {
	Total uint64 `json:"total"`

	Data interface{} `json:"data"`
}
