package user

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AddHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(account.UserBean)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 机构管理员
	agencyId, err := modules.GetAgencyId2(ctx)
	if err != nil {
		logger.Info("GetAgencyId2 err -> ", err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	// 查询是否已经存在的账号
	user, err := account.GetUserByEmail(models.DB(), ctx, req.Email)
	if err != nil {
		logger.Info("GetUserByEmail sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if user != nil {
		logger.Warn(conf.RecordAlreadyExist.String())
		modules.BaseError(ctx, conf.RecordAlreadyExist)
		return
	}

	// 生成新账号
	req.ID = 0
	req.AgencyId = agencyId
	req.Active = true
	if req.AgencyId != 0 { // 普通机构人员创建的用户只允许普通用户，不允许创建管理员账户
		req.Role = string(conf.RoleUser)
	}

	err = models.CreateBaseRecord(req)

	if err != nil {
		logger.Info("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
