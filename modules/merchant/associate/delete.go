package associate

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/merchant"
	"tpayment/modules"
	merchantModule "tpayment/modules/merchant"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func DeleteHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(modules.BaseIDRequest)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 查询是否已经存在的账号
	associateBean, err := merchant.GetUserMerchantAssociateById(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Info("GetUserById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if associateBean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}

	// 判断权限
	err = merchantModule.CheckPermission(ctx, associateBean.MerchantId)
	if err != nil {
		logger.Warn(err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return err
	}

	err = models.DeleteBaseRecord(associateBean)

	if err != nil {
		logger.Info("DeleteBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
