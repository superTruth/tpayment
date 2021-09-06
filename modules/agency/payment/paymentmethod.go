package payment

import (
	"tpayment/conf"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func QueryPaymentMethodsHandle(ctx *gin.Context) {
	logger := tlog.GetGoroutineLogger()

	paymentTypes, err := agency.GetPaymentMethods()
	if err != nil {
		logger.Info("GetPaymentMethods sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	ret := &modules.BaseQueryResponse{
		Total: uint64(len(paymentTypes)),
		Data:  paymentTypes,
	}

	modules.BaseSuccess(ctx, ret)
}
