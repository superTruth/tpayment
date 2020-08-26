package app

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func UpdateHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(tms.App)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 获取机构ID，系统管理员为0
	agencyId, err := modules.GetAgencyId2(ctx)
	if err != nil {
		logger.Warn("GetAgencyId2->", err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	// 查询是否已经存在的账号
	bean, err := tms.GetAppByID(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Info("GetDeviceByID sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if bean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	// 无权限删除
	if agencyId != 0 && agencyId != bean.AgencyId {
		logger.Warn("current agency id is:", bean.AgencyId, ", your id:", agencyId)
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
