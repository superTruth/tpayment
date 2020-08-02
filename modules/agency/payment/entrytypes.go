package payment

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/pkg/tlog"

	"github.com/labstack/echo"
)

func QueryEntryTypesHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	paymentTypes, err := agency.GetEntryTypes(models.DB(), ctx)
	if err != nil {
		logger.Info("GetEntryTypes sql error->", err.Error())
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
