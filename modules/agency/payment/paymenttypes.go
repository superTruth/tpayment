package payment

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func QueryPaymentTypesHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	paymentTypes, err := agency.GetPaymentTypes(models.DB(), ctx)
	if err != nil {
		logger.Info("GetPaymentTypes sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	ret := &modules.BaseQueryResponse{
		Total: uint64(len(paymentTypes)),
		Data:  paymentTypes,
	}

	modules.BaseSuccess(ctx, ret)
}
