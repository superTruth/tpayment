package modules

import (
	"errors"
	"tpayment/conf"
	"tpayment/models/account"
	"tpayment/models/agency"

	"github.com/gin-gonic/gin"
)

func GetAgencyId(ctx *gin.Context, reqAgencyId uint) (uint, error) {
	var agencyId uint
	var userBean *account.UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*account.UserBean)
	} else {
		return 0, errors.New("can't get user")
	}

	if userBean.Role == string(conf.RoleAdmin) {
		if reqAgencyId == 0 {
			return 0, errors.New("Admin user must contain agency id->")
		}
		agencyId = reqAgencyId
	} else {
		var agencys []*agency.Agency
		agencysTmp, ok := ctx.Get(conf.ContextTagAgency)
		if ok {
			agencys = agencysTmp.([]*agency.Agency)
		} else {
			return 0, errors.New("can't get agency")
		}

		agencyId = agencys[0].ID
	}

	return agencyId, nil
}

func GetAgencyId2(ctx *gin.Context) (uint, error) {
	var agencyId uint
	var userBean *account.UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*account.UserBean)
	} else {
		return 0, errors.New("can't get user")
	}

	if userBean.Role == string(conf.RoleAdmin) {
		agencyId = 0
	} else {
		var agencys []*agency.Agency
		agencysTmp, ok := ctx.Get(conf.ContextTagAgency)
		if ok {
			agencys = agencysTmp.([]*agency.Agency)
		} else {
			return 0, errors.New("can't get agency")
		}

		if len(agencys) == 0 {
			return 0, errors.New("no permission")
		}

		agencyId = agencys[0].ID
	}

	return agencyId, nil
}

func IsAgencyAdmin(ctx *gin.Context) *agency.Agency {
	var agencys []*agency.Agency
	agencysTmp, ok := ctx.Get(conf.ContextTagAgency)
	if ok {
		agencys = agencysTmp.([]*agency.Agency)
	} else {
		return nil
	}

	if len(agencys) == 0 {
		return nil
	}

	return agencys[0]
}

func IsAdmin(ctx *gin.Context) *account.UserBean {
	var userBean *account.UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*account.UserBean)
	} else {
		return nil
	}

	if userBean.Role != string(conf.RoleAdmin) {
		return nil
	}

	return userBean
}
