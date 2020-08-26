package tms

import (
	"errors"
	"tpayment/conf"
	"tpayment/models/account"
	"tpayment/models/agency"
	"tpayment/models/tms"
	"tpayment/modules"

	"github.com/gin-gonic/gin"
)

func CheckPermission(ctx *gin.Context, deviceInfo *tms.DeviceInfo) error {

	var userBean *account.UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*account.UserBean)
	} else {
		return errors.New("can't get user")
	}

	// 系统管理员可以操作一切
	if userBean.Role == string(conf.RoleAdmin) {
		return nil
	}

	var agencys []*agency.Agency
	agencysTmp, ok := ctx.Get(conf.ContextTagAgency)
	if ok {
		agencys = agencysTmp.([]*agency.Agency)
	} else {
		modules.BaseError(ctx, conf.UnknownError)
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
