package agency

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"tpayment/models"
)

type UserAgencyAssociate struct {
	gorm.Model

	AgencyId uint `json:"agency_id" gorm:"column:agency_id"`
	UserId   uint `json:"user_id" gorm:"column:user_id"`
}

func (UserAgencyAssociate) TableName() string {
	return "agency_user_associate"
}

func GetAssociateById(db *models.MyDB, ctx echo.Context, id uint) (*UserAgencyAssociate, error) {
	ret := new(UserAgencyAssociate)

	err := db.Model(&UserAgencyAssociate{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetAssociateByUserId(db *models.MyDB, ctx echo.Context, userId uint) (*UserAgencyAssociate, error) {
	ret := new(UserAgencyAssociate)

	err := db.Model(&UserAgencyAssociate{}).Where("user_id=?", userId).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

type AssociateAgencyUserBean struct {
	gorm.Model

	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func QueryUsersByAgencyId(db *models.MyDB, ctx echo.Context, agencyId, offset, limit uint, filters map[string]string) (uint, []AssociateAgencyUserBean, error) {
	filterTmp := make(map[string]interface{})

	for k, v := range filters {
		filterTmp[k] = v
	}

	// conditions
	tmpDb := db.Table("user").Where(filterTmp)
	tmpDb = tmpDb.Joins("JOIN agency_user_associate ass ON ass.agency_id = ? AND ass.user_id = user.id AND ass.deleted_at IS NULL", agencyId)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []AssociateAgencyUserBean
	if err = tmpDb.Offset(offset).Limit(limit).Select("user.id as id, ass.created_at as created_at, ass.updated_at as updated_at, user.name as name, user.email as email").Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
