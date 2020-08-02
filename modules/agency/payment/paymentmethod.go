package payment

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/pkg/tlog"

	"github.com/labstack/echo"
)

func QueryPaymentMethodsHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	paymentTypes, err := agency.GetPaymentMethods(models.DB(), ctx)
	if err != nil {
		logger.Info("GetPaymentMethods sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	ret := &modules.BaseQueryResponse{
		Total: uint(len(paymentTypes)),
		Data:  paymentTypes,
	}

	modules.BaseSuccess(ctx, ret)

	return nil
}
