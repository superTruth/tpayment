package user

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/models/agency"
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

	// 机构管理员
	agencyId := uint64(0)
	var userBean *account.UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*account.UserBean)
	} else {
		modules.BaseError(ctx, conf.UnknownError)
		return
	}

	var agencys []*agency.Agency
	agencysTmp, ok := ctx.Get(conf.ContextTagAgency)
	if ok {
		agencys = agencysTmp.([]*agency.Agency)
	} else {
		modules.BaseError(ctx, conf.UnknownError)
		return
	}

	if userBean.Role != string(conf.RoleAdmin) { // 管理员，不需要过滤机构
		if len(agencys) == 0 {
			logger.Info("not admin and not agency")
			modules.BaseError(ctx, conf.ParameterError)
			return
		}
		agencyId = agencys[0].ID
	}

	total, dataRet, err := account.QueryUserRecord(models.DB(), ctx, req.Offset, req.Limit, agencyId, req.Filters)
	if err != nil {
		logger.Info("QueryBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	ret := &modules.BaseQueryResponse{
		Total: total,
		Data:  dataRet,
	}

	modules.BaseSuccess(ctx, ret)
}
