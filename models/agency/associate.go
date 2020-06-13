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

func GetAssociateByUserId(db *models.MyDB, ctx echo.Context, userID uint) (*UserAgencyAssociate, error) {
	ret := new(UserAgencyAssociate)

	err := db.Model(&UserAgencyAssociate{}).Where("user_id=?", userID).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}
