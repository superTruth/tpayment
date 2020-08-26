package associate

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func DeleteHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(modules.BaseIDRequest)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	if req.ID == 0 {
		logger.Warn("ParameterError")
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 查询是否已经存在的账号
	associateBean, err := agency.GetAssociateById(models.DB(), ctx, req.ID)
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

	err = models.DeleteBaseRecord(associateBean)

	if err != nil {
		logger.Info("DeleteBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
