package modules

import (
	"errors"
	"tpayment/conf"
	"tpayment/models/account"
	"tpayment/models/agency"

	"github.com/labstack/echo"
)

func GetAgencyId(ctx echo.Context, reqAgencyId uint) (uint, error) {
	var agencyId uint
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
	if userBean.Role == string(conf.RoleAdmin) {
		if reqAgencyId == 0 {
			return 0, errors.New("Admin user must contain agency id->")
		}
		agencyId = reqAgencyId
	} else {
		agencys := ctx.Get(conf.ContextTagAgency).([]*agency.Agency)
		agencyId = agencys[0].ID
	}

	return agencyId, nil
}

func GetAgencyId2(ctx echo.Context) (uint, error) {
	var agencyId uint
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
	if userBean.Role == string(conf.RoleAdmin) {
		agencyId = 0
	} else {
		agencys, ok := ctx.Get(conf.ContextTagAgency).([]*agency.Agency)
		if !ok || len(agencys) == 0 {
			return 0, errors.New("no permission")
		}

		agencyId = agencys[0].ID
	}

	return agencyId, nil
}

func IsAgencyAdmin(ctx echo.Context) *agency.Agency {
	agencys, ok := ctx.Get(conf.ContextTagAgency).([]*agency.Agency)
	if !ok {
		return nil
	}

	if len(agencys) == 0 {
		return nil
	}

	return agencys[0]
}

func IsAdmin(ctx echo.Context) *account.UserBean {
	userBean, ok := ctx.Get(conf.ContextTagUser).(*account.UserBean)

	if !ok {
		return nil
	}

	if userBean.Role != string(conf.RoleAdmin) {
		return nil
	}

	return userBean
}
