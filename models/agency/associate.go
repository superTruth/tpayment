package agency

import (
	"github.com/jinzhu/gorm"
)

type UserAgencyAssociate struct {
	gorm.Model

	AgencyId uint `json:"agency_id" gorm:"column:agency_id"`
	UserId   uint `json:"user_id" gorm:"column:user_id"`
}

func (UserAgencyAssociate) TableName() string {
	return "agency_user_associate"
}

