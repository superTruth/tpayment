package devicetag

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func DeleteHandle(ctx *gin.Context) {
	logger := tlog.GetGoroutineLogger()

	req := new(modules.BaseIDRequest)

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
	bean, err := tms.GetDeviceTagByID(req.ID)
	if err != nil {
		logger.Info("GetAppInDeviceByID sql error->", err.Error())
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

	// 数据正在使用，无法删除
	isUsing, err := tms.IsTagUsing(req.ID)
	if err != nil {
		logger.Info("IsTagUsing sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if isUsing {
		logger.Info("tag is using->", req.ID)
		modules.BaseError(ctx, conf.DataIsUsing)
		return
	}

	err = models.DeleteBaseRecord(bean)

	if err != nil {
		logger.Info("DeleteBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
