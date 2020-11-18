package user

import (
	"tpayment/conf"
	"tpayment/internal/basekey"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func UpdateHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(account.UserBean)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 密码进行hash
	if req.Pwd != "" {
		req.Pwd = basekey.Hash([]byte(req.Pwd))
	}

	// 查询是否已经存在的账号
	user, err := account.GetUserById(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Info("GetUserById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if user == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordAlreadyExist)
		return
	}

	// 管理员用户，可以更改一切数据
	if modules.IsAdmin(ctx) != nil {

	} else {
		// 机构用户，可以更改自己创建的账号数据和自己的数据，不能更改角色
		// 普通用户，只能改自己的数据，不能更改角色
		req.Role = ""

		// 判断目标数据是否就是自己数据
		var currentUserBean *account.UserBean
		currentUserBeanTmp, ok := ctx.Get(conf.ContextTagUser)
		if ok {
			currentUserBean = currentUserBeanTmp.(*account.UserBean)
		} else {
			modules.BaseError(ctx, conf.UnknownError)
			return
		}

		if currentUserBean.ID == user.ID { // 是自己数据的情况，直接运行修改

		} else {
			agency := modules.IsAgencyAdmin(ctx)
			if agency == nil { // 不是机构账号，也不是自己数据，不允许修改
				logger.Warn("can't update other user account: your account id->", currentUserBean.ID,
					"  dest account id->", user.ID)
				modules.BaseError(ctx, conf.NoPermission)
				return
			}

			// 机构账号修改的不是自己创建的数据
			if agency.ID != user.AgencyId {
				logger.Warn("the user is not belong your agency->", user.ID)
				modules.BaseError(ctx, conf.NoPermission)
				return
			}
		}

	}
	//
	//// 更改权限
	//if modules.IsAdmin(ctx) == nil { // 如果不是系统管理员，则不允许更改角色
	//	req.Role = ""
	//} else { // 查看是否是机构管理员，如果不是，只能更新自己的数据
	//	agency := modules.IsAgencyAdmin(ctx)
	//	if agency == nil {
	//		var currentUserBean *account.UserBean
	//		currentUserBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	//		if ok {
	//			currentUserBean = currentUserBeanTmp.(*account.UserBean)
	//		} else {
	//			modules.BaseError(ctx, conf.UnknownError)
	//			return
	//		}
	//
	//		if currentUserBean.ID != user.ID {
	//			logger.Warn("can't update other user account: your account id->", currentUserBean.ID,
	//				"  dest account id->", user.ID)
	//			modules.BaseError(ctx, conf.NoPermission)
	//			return
	//		}
	//	} else {
	//		// 修改用户不属于此机构用户
	//		if user.AgencyId != agency.ID {
	//			logger.Warn("can't update other user account: your agency id->", agency.ID,
	//				"  dest account id->", user.ID)
	//			modules.BaseError(ctx, conf.NoPermission)
	//			return
	//		}
	//	}
	//}

	// 生成新账号
	req.AgencyId = 0
	req.Active = user.Active
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)
}
