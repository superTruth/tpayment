package tms

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/models/agency"
)

type App struct {
	gorm.Model

	AgencyId    uint   `gorm:"column:agency_id"`
	Name        string `gorm:"column:name"`
	PackageId   string `gorm:"column:package_id"`
	Description string `gorm:"column:description"`
}

func (App) TableName() string {
	return "mdm2_apps"
}

// 根据device ID获取设备信息
func GetAppByID(db *models.MyDB, ctx echo.Context, id uint) (*App, error) {

	ret := new(App)

	err := db.Model(&App{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryAppRecord(db *models.MyDB, ctx echo.Context, offset, limit uint, filters map[string]string) (uint, []App, error) {
	filterTmp := make(map[string]interface{})
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
	agencys := ctx.Get(conf.ContextTagAgency).([]agency.Agency)

	for k, v := range filters {
		filterTmp[k] = v
	}

	if userBean.Role != string(conf.RoleAdmin) { // 管理员，不需要过滤机构
		if len(agencys) == 0 {
			return 0, nil, errors.New("user not agency admin")
		}
		filterTmp["agency_id"] = agencys[0].ID
	}

	// conditions
	tmpDb := db.Table("mdm2_apps").Where(filterTmp)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []App
	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
