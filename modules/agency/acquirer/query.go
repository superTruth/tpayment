package acquirer

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
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

	// 管理员必须要传入agency id
	agencyId, err := modules.GetAgencyId(ctx, req.AgencyId)
	if err != nil {
		logger.Warn("GetAgencyId no permission->", err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return err
	}

	if req.Limit > conf.MaxQueryCount { // 一次性不能搜索太多数据
		req.Limit = conf.MaxQueryCount
	}

	total, dataRet, err := agency.QueryAcquirerRecord(models.DB(), ctx, agencyId, req.Offset, req.Limit, req.Filters)
	if err != nil {
		logger.Info("QueryBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	ret := &modules.BaseQueryResponse{
		Total: total,
		Data:  dataRet,
	}

	modules.BaseSuccess(ctx, ret)

	return nil
}
