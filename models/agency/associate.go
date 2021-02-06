package agency

import (
	"tpayment/models"
	"tpayment/models/account"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserAgencyAssociate struct {
	models.BaseModel

	AgencyId uint64 `gorm:"column:agency_id" json:"agency_id"`
	UserId   uint64 `gorm:"column:user_id" json:"user_id"`
}

func (UserAgencyAssociate) TableName() string {
	return "agency_user_associate"
}

func GetAssociateById(db *models.MyDB, ctx *gin.Context, id uint64) (*UserAgencyAssociate, error) {
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

func GetAssociateByUserId(db *models.MyDB, ctx *gin.Context, userId uint64) (*UserAgencyAssociate, error) {
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
	models.BaseModel

	Email string `gorm:"column:email" json:"email"`
	Name  string `gorm:"column:name" json:"name"`
}

func QueryUsersByAgencyId(db *models.MyDB, ctx *gin.Context, agencyId, offset, limit uint64) (uint64, []*AssociateAgencyUserBean, error) {
	// conditions
	tmpDb := db.Table("user").Model(&account.UserBean{})
	tmpDb = tmpDb.Joins("JOIN agency_user_associate ass ON ass.agency_id = ? AND ass.user_id = user.id AND ass.deleted_at IS NULL", agencyId)

	// 统计总数
	var total uint64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*AssociateAgencyUserBean
	if err = tmpDb.Order("id desc").Offset(offset).Limit(limit).Select(
		"ass.id as id, ass.created_at as created_at, " +
			"ass.updated_at as updated_at, user.name as name, " +
			"user.email as email").Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
