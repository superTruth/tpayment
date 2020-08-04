package agency

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Agency struct {
	models.BaseModel

	Name string `gorm:"column:name" json:"name"`
	Tel  string `gorm:"column:tel" json:"tel"`
	Addr string `gorm:"column:addr" json:"addr"`
}

func (Agency) TableName() string {
	return "agency"
}

func QueryAgencyRecord(db *models.MyDB, ctx echo.Context, offset, limit uint, filters map[string]string) (uint, []*Agency, error) {
	var ret []*Agency

	equalData := make(map[string]string)
	sqlCondition := models.CombQueryCondition(equalData, filters)

	tmpDB := db.Model(&Agency{}).Where(sqlCondition)

	// 非系统管理员，只查看跟他有关的机构
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
	if userBean.Role != string(conf.RoleAdmin) {
		tmpDB = tmpDB.Joins("JOIN agency_user_associate ass ON ass.user_id = ? AND ass.deleted_at IS NULL", userBean.ID)
	}

	// 统计总数
	var total uint = 0
	err := tmpDB.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	// 查询记录
	err = tmpDB.Offset(offset).Limit(limit).Find(&ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return 0, ret, nil
		}
		return 0, nil, err
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
