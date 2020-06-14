package merchant

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
)

type Merchant struct {
	gorm.Model

	AgencyId uint   `json:"agency_id" gorm:"column:agency_id"`
	Name     string `json:"name" gorm:"column:name"`
	Tel      string `json:"tel"  gorm:"column:tel"`
	Addr     string `json:"addr" gorm:"column:addr"`
}

func (Merchant) TableName() string {
	return "merchant"
}

func GetMerchantById(id uint) (*Merchant, error) {
	ret := new(Merchant)

	err := models.DB().Model(&Merchant{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

//func QueryMerchantRecord(offset, limit uint,filters map[string]string) (uint, []Merchant, error) {
//	filterTmp := make(map[string]interface{})
//
//	for k,v := range filters {
//		filterTmp[k] = v
//	}
//
//	// 统计总数
//	var total uint = 0
//	err := models.DB().Model(&Merchant{}).Where(filterTmp).Count(&total).Error
//	if err != nil {
//		return 0, nil, err
//	}
//
//	// 查询记录
//	var ret []Merchant
//
//	err = models.DB().Model(&Merchant{}).Where(filterTmp).Offset(offset).Limit(limit).Find(&ret).Error
//
//	if err != nil {
//		if gorm.ErrRecordNotFound == err { // 没有记录
//			return 0, ret, nil
//		}
//		return 0, nil, err
//	}
//
//	return total, ret, nil
//}

func QueryMerchantRecord(db *models.MyDB, ctx echo.Context, agencyId, offset, limit uint, filters map[string]string) (uint, []Merchant, error) {
	filterTmp := make(map[string]interface{})
	userBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)

	for k, v := range filters {
		filterTmp[k] = v
	}

	if agencyId != 0 {
		filterTmp["agency_id"] = agencyId
	}

	// conditions
	tmpDb := db.Table("merchant").Where(filterTmp)
	if userBean.Role != string(conf.RoleAdmin) {  // 管理员账户可以忽略这个选项
		tmpDb = tmpDb.Joins("JOIN merchant_user_associate ass ON ass.merchant_id = merchant.id AND ass.user_id = ? AND ass.deleted_at IS NULL", userBean.ID)
	}

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []Merchant
	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

func GetMerchantsById(ids []uint) ([]Merchant, error) {
	var ret []Merchant

	err := models.DB().Model(&Merchant{}).Where(ids).First(&ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}
