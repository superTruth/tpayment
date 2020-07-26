package appinbatchupdate

import (
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"
)

// TODO 未完成
func AddHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(tms.AppInDevice)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	//errorCode := SmartAddAppInDevice(ctx, req)
	//if errorCode != conf.SUCCESS {
	//	modules.BaseError(ctx, errorCode)
	//	return errors.New(errorCode.String())
	//}

	modules.BaseSuccess(ctx, nil)

	return nil
}
