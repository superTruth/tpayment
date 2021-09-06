package associate

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

func AddHandle(ctx *gin.Context) {
	logger := tlog.GetGoroutineLogger()

	req := new(agency.UserAgencyAssociate)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	if req.AgencyId == 0 || req.UserId == 0 {
		logger.Warn("ParameterError")
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 查询是否存在这2个ID
	userBean, err := account.GetUserById(req.UserId)
	if err != nil {
		logger.Info("GetUserById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if userBean == nil {
		logger.Warn("User Not Exist")
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	agencyBean, err := agency.GetAgencyById(req.AgencyId)
	if err != nil {
		logger.Info("GetAssociateById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if agencyBean == nil {
		logger.Warn("Agency Not Exist")
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 一个用户只可以关联一个agency
	associateBean, err := agency.GetAssociateByUserId(req.UserId)
	if err != nil {
		logger.Info("GetAssociateByUserId sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if associateBean != nil {
		logger.Warn(conf.UserCanOnlyInOneAgency.String())
		modules.BaseError(ctx, conf.UserCanOnlyInOneAgency)
		return
	}

	// 如果是系统管理员，则无法关联
	if userBean.Role != string(conf.RoleUser) { // 只有普通用户可以关联
		logger.Info(conf.AdminCantAssociate.String())
		modules.BaseError(ctx, conf.AdminCantAssociate)
		return
	}

	err = models.CreateBaseRecord(req)
	if err != nil {
		logger.Info("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
