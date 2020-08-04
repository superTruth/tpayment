package device

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func QueryHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(modules.BaseQueryRequest)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	if req.Limit > conf.MaxQueryCount { // 一次性不能搜索太多数据
		req.Limit = conf.MaxQueryCount
	}

	total, dataRet, err := tms.QueryDeviceRecord(models.DB(), ctx, req.Offset, req.Limit, req.Filters)
	if err != nil {
		logger.Info("QueryBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	// 查询出设备里面对应的所有的tag
	for i := 0; i < len(dataRet); i++ {
		tags, err := tms.QueryTagsInDevice(models.DB(), ctx, dataRet[i])
		if err != nil {
			logger.Info("QueryTagsInDevice sql error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return err
		}
		dataRet[i].Tags = &tags
	}

	ret := &modules.BaseQueryResponse{
		Total: total,
		Data:  dataRet,
	}

	modules.BaseSuccess(ctx, ret)

	return nil
}
