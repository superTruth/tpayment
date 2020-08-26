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

func AddHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(tms.DeviceTag)

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

	req.AgencyId = agencyId
	req.ID = 0
	err = models.CreateBaseRecord(req)

	if err != nil {
		logger.Error("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)

}
