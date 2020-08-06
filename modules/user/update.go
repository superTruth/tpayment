package user

import (
	"fmt"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func UpdateHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)
	fmt.Println("user UpdateHandle")

	req := new(account.UserBean)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 查询是否已经存在的账号
	user, err := account.GetUserById(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Info("GetUserById sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if user == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordAlreadyExist)
		return err
	}

	// 更改权限
	if modules.IsAdmin(ctx) == nil { // 如果不是系统管理员，则不允许更改角色
		req.Role = ""
	} else { // 查看是否是机构管理员，如果不是，只能更新自己的数据
		agency := modules.IsAgencyAdmin(ctx)
		if agency == nil {
			currentUserBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
			if currentUserBean.ID != user.ID {
				logger.Warn("can't update other user account: your account id->", currentUserBean.ID,
					"  dest account id->", user.ID)
				modules.BaseError(ctx, conf.NoPermission)
				return err
			}
		} else {
			// 修改用户不属于此机构用户
			if user.AgencyId != agency.ID {
				logger.Warn("can't update other user account: your agency id->", agency.ID,
					"  dest account id->", user.ID)
				modules.BaseError(ctx, conf.NoPermission)
				return err
			}
		}
	}

	// 生成新账号
	req.AgencyId = 0
	req.Active = user.Active
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
