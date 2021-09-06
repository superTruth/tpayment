package devicetag

import (
	"tpayment/conf"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func DeviceInTagHandle(ctx *gin.Context) {
	logger := tlog.GetGoroutineLogger()

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

	total, dataRet, err := tms.DeviceTagDao.QueryDeviceInTag(req.TagId, req.Offset, req.Limit, req.Filters)
	if err != nil {
		logger.Info("QueryDeviceInTag sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	ret := &modules.BaseQueryResponse{
		Total: total,
		Data:  dataRet,
	}

	modules.BaseSuccess(ctx, ret)
}
