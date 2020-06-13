package agency

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
)

type Agency struct {
	gorm.Model

	Name string `json:"name" gorm:"column:name"`
	Tel  string `json:"tel"  gorm:"column:tel"`
	Addr string `json:"addr" gorm:"column:addr"`
}

func (Agency) TableName() string {
	return "agency"
}

func QueryAgencyRecord(db *models.MyDB, ctx echo.Context, offset, limit uint, filters map[string]string) (uint, []Agency, error) {
	filterTmp := make(map[string]interface{})
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)

	for k, v := range filters {
		filterTmp[k] = v
	}

	// conditions
	tmpDb := db.Table("agency").Where(filterTmp)
	if userBean.Role != string(conf.RoleAdmin) {
		tmpDb = tmpDb.Joins("JOIN agency_user_associate ass ON ass.agency_id = agency.id AND ass.user_id = ? AND ass.deleted_at IS NULL", userBean.ID)
	}

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []Agency
	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

func GetAgencyById(db *models.MyDB, ctx echo.Context, id uint) (*Agency, error) {
	ret := new(Agency)

	err := db.Model(&Agency{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}
