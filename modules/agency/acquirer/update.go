package acquirer

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

func UpdateHandle(ctx *gin.Context) {
	logger := tlog.GetGoroutineLogger()

	req := new(agency.Acquirer)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 查询是否已经存在的账号
	acquirerBean, err := agency.GetAcquirerById(req.ID)
	if err != nil {
		logger.Info("GetMerchantById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if acquirerBean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	// 判断当前agency是否有权限
	var userBean *account.UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*account.UserBean)
	} else {
		modules.BaseError(ctx, conf.UnknownError)
		return
	}
	if userBean.Role != string(conf.RoleAdmin) {
		var agencys []*agency.Agency
		agencysTmp, ok := ctx.Get(conf.ContextTagAgency)
		if ok {
			agencys = agencysTmp.([]*agency.Agency)
		} else {
			modules.BaseError(ctx, conf.UnknownError)
			return
		}

		if agencys[0].ID != acquirerBean.AgencyId {
			logger.Warn("this acquirer is not belong to the agency")
			modules.BaseError(ctx, conf.NoPermission)
			return
		}
	}

	// 生成新账号
	req.AgencyId = 0 // 不允许更新agency id
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
