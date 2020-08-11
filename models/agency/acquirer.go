package agency

import (
	"strconv"
	"tpayment/models"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Acquirer struct {
	models.BaseModel

	Name          string `json:"name" gorm:"column:name"`
	Addition      string `json:"addition"  gorm:"column:addition"`
	ConfigFileUrl string `json:"config_file_url" gorm:"column:config_file_url"`
	AgencyId      uint   `json:"agency_id"  gorm:"column:agency_id"`
}

func (Acquirer) TableName() string {
	return "acquirer"
}

func GetAcquirerById(id uint) (*Acquirer, error) {
	ret := new(Acquirer)

	err := models.DB().Model(&Acquirer{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryAcquirerRecord(db *models.MyDB, ctx echo.Context, agencyId, offset, limit uint, filters map[string]string) (uint, []*Acquirer, error) {
	equalData := make(map[string]string)
	if agencyId != 0 {
		equalData["agency_id"] = strconv.FormatUint(uint64(agencyId), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Model(&Acquirer{}).Where(sqlCondition)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*Acquirer
	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
