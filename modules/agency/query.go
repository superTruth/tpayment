package agency

import (
	"tpayment/conf"
	"tpayment/models/account"
	"tpayment/models/agency"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func QueryHandle(ctx *gin.Context) {
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

	total, dataRet, err := agency.QueryAgencyRecord(ctx, req.Offset, req.Limit, req.Filters)
	if err != nil {
		logger.Info("QueryBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	// 获取agency的权限
	agencyId, err := modules.GetAgencyId2(ctx)
	if err != nil {
		logger.Warn("GetAgencyId no permission->", err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	if agencyId != 0 {
		var userBean *account.UserBean
		userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
		if ok {
			userBean = userBeanTmp.(*account.UserBean)
		} else {
			return
		}

		for i := 0; i < len(dataRet); i++ {
			roleBean, err := agency.UserAgencyAssociateDao.GetByAgencyUserID(dataRet[i].ID, userBean.ID)
			if err != nil {
				logger.Errorf("GetByAgencyUserID sql error->", err.Error())
				modules.BaseError(ctx, conf.DBError)
				return
			}
			if roleBean.Role == "" {
				roleBean.Role = string(conf.MerchantManager)
			}
			dataRet[i].Role = roleBean.Role
		}
	}

	ret := &modules.BaseQueryResponse{
		Total: total,
		Data:  dataRet,
	}

	modules.BaseSuccess(ctx, ret)
}
