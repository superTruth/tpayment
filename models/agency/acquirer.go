package agency

import (
	"strconv"
	"tpayment/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var AcquirerDao = &Acquirer{}

type Acquirer struct {
	models.BaseModel

	Name               string `json:"name" gorm:"column:name"`
	ImplName           string `json:"impl_name" gorm:"column:impl_name"`
	Addition           string `json:"addition"  gorm:"column:addition"`
	ConfigFileUrl      string `json:"config_file_url" gorm:"column:config_file_url"`
	AgencyId           uint64 `json:"agency_id"  gorm:"column:agency_id"`
	AutoSettlementTime string `json:"auto_settlement_time" gorm:"column:auto_settlement_time"`
}

func (Acquirer) TableName() string {
	return "agency_acquirer"
}

func (acq *Acquirer) Get(id uint64) (*Acquirer, error) {
	ret := new(Acquirer)

	err := models.DB().Model(acq).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetAcquirerById(id uint64) (*Acquirer, error) {
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

func QueryAcquirerRecord(db *models.MyDB, ctx *gin.Context, agencyId, offset, limit uint64, filters map[string]string) (uint64, []*Acquirer, error) {
	equalData := make(map[string]string)
	if agencyId != 0 {
		equalData["agency_id"] = strconv.FormatUint(uint64(agencyId), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Model(&Acquirer{}).Where(sqlCondition)

	// 统计总数
	var total uint64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*Acquirer
	if err = tmpDb.Order("id desc").Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

// 查找需要结算的收单
func (acq *Acquirer) GetNeedSettlement(hour string) ([]*Acquirer, error) {
	var ret []*Acquirer

	err := models.DB().Model(acq).Where("auto_settlement_time LIKE ?%", hour).Find(&ret).Error

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (acq *Acquirer) GetByName(agencyID uint64, name string) (*Acquirer, error) {
	ret := &Acquirer{}
	err := models.DB().Model(ret).
		Where("agency_id = ? and name = ?", agencyID, name).
		First(&ret).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}
