package tms

import (
	"errors"
	"tpayment/conf"
	"tpayment/models/account"
	"tpayment/models/agency"
	"tpayment/models/tms"

	"github.com/labstack/echo"
)

func CheckPermission(ctx echo.Context, deviceInfo *tms.DeviceInfo) error {
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)

	// 系统管理员可以操作一切
	if userBean.Role == string(conf.RoleAdmin) {
		return nil
	}

	agencys, ok := ctx.Get(conf.ContextTagAgency).([]*agency.Agency)
	if !ok {
		return errors.New("get agency fail")
	}

	if len(agencys) == 0 {
		return errors.New("not agency")
	}

	if agencys[0].ID != deviceInfo.AgencyId {
		return errors.New("no permission")
	}

	return nil
}
