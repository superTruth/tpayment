package merchant

import (
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/merchant"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"
)

func UpdateHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(merchant.Merchant)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 查询是否已经存在的账号
	merchantBean, err := merchant.GetMerchantById(req.ID)
	if err != nil {
		logger.Info("GetMerchantById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if merchantBean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}

	// 生成新账号
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	//
	ret := &modules.BaseResponse{
		ErrorCode:    conf.SUCCESS,
	}

	modules.BaseSuccess(ctx, ret)

	return nil
}

