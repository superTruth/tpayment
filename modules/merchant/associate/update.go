package associate

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/merchant"
	"tpayment/modules"
	merchantModule "tpayment/modules/merchant"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func UpdateHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(merchant.UserMerchantAssociate)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 查询是否已经存在的账号
	associateBean, err := merchant.GetUserMerchantAssociateById(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Info("GetUserById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if associateBean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	// 判断权限
	err = merchantModule.CheckPermission(ctx, associateBean.MerchantId)
	if err != nil {
		logger.Warn(err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	// 更新数据
	req.MerchantId = 0 // 不允许更新这2个数据
	req.UserId = 0
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
