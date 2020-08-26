package merchant

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/merchant"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func UpdateHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(merchant.Merchant)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 查询是否已经存在的账号
	merchantBean, err := merchant.GetMerchantById(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Info("GetMerchantById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if merchantBean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	// 判断是否是属于自己机构的商户
	_, err = modules.GetAgencyId(ctx, merchantBean.AgencyId)
	if err != nil {
		logger.Warn(err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	// 生成新账号
	req.AgencyId = 0
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
