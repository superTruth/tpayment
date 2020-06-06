package agency

import (
	"github.com/jinzhu/gorm"
	"tpayment/models"
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

func QueryAgencyRecord(userId, offset, limit uint, filters map[string]string) (uint, []Agency, error) {
	filterTmp := make(map[string]interface{})

	for k, v := range filters {
		filterTmp[k] = v
	}

	// conditions
	tmpDb := models.DB().Table("agency").Where(filterTmp)
	if userId != 0 {
		tmpDb = tmpDb.Joins("JOIN agency_user_associate ass ON ass.agency_id = agency.id AND ass.user_id = ? AND ass.deleted_at IS NULL", userId)
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

func GetAgencyById(id uint) (*Agency, error) {
	ret := new(Agency)

	err := models.DB().Model(&Agency{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}
