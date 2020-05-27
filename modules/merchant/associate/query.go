package associate

import (
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models/account"
	"tpayment/models/merchant"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"
)

func QueryHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)

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

	var userId uint = 0
	if userBean.Role != string(conf.RoleAdmin) {  // 管理员用户可以搜索所有商户
		userId = userBean.ID
	}

	total, dataRet, err := merchant.QueryMerchantRecord(userId, req.Offset, req.Limit, req.Filters)
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

func QueryUserInMerchantHandle(ctx echo.Context) error {
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

	total, dataRet, err := merchant.QueryUsersByMerchantId(req.MerchantId, req.Offset, req.Limit, req.Filters)
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
