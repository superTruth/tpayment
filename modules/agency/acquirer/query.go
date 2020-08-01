package acquirer

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
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
	var agencyId uint
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
	if userBean.Role == string(conf.RoleAdmin) {
		if req.AgencyId == 0 {
			logger.Warn("Admin user must contain agency id->")
			modules.BaseError(ctx, conf.ParameterError)
			return err
		}
		agencyId = req.AgencyId
	} else {
		agencys := ctx.Get(conf.ContextTagAgency).([]*agency.Agency)
		agencyId = agencys[0].ID
		req.AgencyId = agencyId
	}

	if req.Limit > conf.MaxQueryCount { // 一次性不能搜索太多数据
		req.Limit = conf.MaxQueryCount
	}

	total, dataRet, err := agency.QueryAcquirerRecord(models.DB(), ctx, req.AgencyId, req.Offset, req.Limit, req.Filters)
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
