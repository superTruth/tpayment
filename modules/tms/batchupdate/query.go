package batchupdate

import (
	"fmt"
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

	total, dataRet, err := tms.QueryBatchUpdateRecord(models.DB(), ctx, req.Offset, req.Limit, req.Filters)
	if err != nil {
		logger.Info("QueryAppInDeviceRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	for i := 0; i < len(dataRet); i++ {
		fmt.Println("Truth device tags->", dataRet[i].Tags)
		dataRet[i].ConfigTags, err = tms.GetDeviceTagByIDs(models.DB(), ctx, dataRet[i].Tags)
		if err != nil {
			logger.Error("GetDeviceTagByIDs fail->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return
		}

		fmt.Println("Truth device models->", dataRet[i].DeviceModels)
		dataRet[i].ConfigModels, err = tms.GetModelByIDs(models.DB(), ctx, dataRet[i].DeviceModels)
		if err != nil {
			logger.Error("GetModelByIDs fail->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return
		}
	}

	ret := &modules.BaseQueryResponse{
		Total: total,
		Data:  dataRet,
	}

	modules.BaseSuccess(ctx, ret)
}
