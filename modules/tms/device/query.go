package device

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func QueryHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(modules.BaseQueryRequest)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	if req.Limit > conf.MaxQueryCount { // 一次性不能搜索太多数据
		req.Limit = conf.MaxQueryCount
	}

	total, dataRet, err := tms.QueryDeviceRecord(models.DB(), ctx, req.Offset, req.Limit, req.Filters)
	if err != nil {
		logger.Info("QueryBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	// 查询出设备里面对应的所有的tag
	for i := 0; i < len(dataRet); i++ {
		tags, err := tms.QueryTagsInDevice(models.DB(), ctx, dataRet[i])
		if err != nil {
			logger.Info("QueryTagsInDevice sql error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return
		}
		dataRet[i].Tags = &tags
	}

	// 查询出所有设备对应的device model
	for i := 0; i < len(dataRet); i++ {
		if dataRet[i].DeviceModel == 0 {
			continue
		}
		deviceModel, err := tms.GetModelByID(models.DB(), ctx, dataRet[i].DeviceModel)
		if err != nil {
			logger.Info("GetModelByID sql error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return
		}
		dataRet[i].DeviceModelName = deviceModel.Name
	}

	ret := &modules.BaseQueryResponse{
		Total: total,
		Data:  dataRet,
	}

	modules.BaseSuccess(ctx, ret)
}
