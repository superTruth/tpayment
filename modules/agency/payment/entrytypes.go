package payment

import (
	"tpayment/conf"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func QueryEntryTypesHandle(ctx *gin.Context) {
	logger := tlog.GetGoroutineLogger()

	paymentTypes, err := agency.GetEntryTypes()
	if err != nil {
		logger.Info("GetEntryTypes sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	ret := &modules.BaseQueryResponse{
		Total: uint64(len(paymentTypes)),
		Data:  paymentTypes,
	}

	modules.BaseSuccess(ctx, ret)
}
