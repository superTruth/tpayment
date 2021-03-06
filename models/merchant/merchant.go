package merchant

import (
	"errors"
	"strconv"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/modules"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Merchant struct {
	models.BaseModel

	AgencyId uint   `json:"agency_id" gorm:"column:agency_id"`
	Name     string `json:"name,omitempty" gorm:"column:name"`
	Tel      string `json:"tel,omitempty"  gorm:"column:tel"`
	Addr     string `json:"addr,omitempty" gorm:"column:addr"`
	Email    string `json:"email,omitempty" gorm:"column:email"`
}

func (Merchant) TableName() string {
	return "merchant"
}

func GetMerchantById(db *models.MyDB, ctx *gin.Context, id uint) (*Merchant, error) {
	ret := new(Merchant)

	err := db.Model(&Merchant{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

// 获取机构下面的商户
func QueryMerchantInAgency(db *models.MyDB, ctx *gin.Context, agencyId, offset, limit uint, filters map[string]string) (uint, []*Merchant, error) {
	equalData := make(map[string]string)
	if agencyId != 0 {
		equalData["agency_id"] = strconv.FormatUint(uint64(agencyId), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Model(&Merchant{}).Where(sqlCondition)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*Merchant
	if err = tmpDb.Order("updated_at desc").Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

// 获取账号相关的商户列表
func QueryMerchantInUser(db *models.MyDB, ctx *gin.Context, offset, limit uint, filters map[string]string) (uint, []*Merchant, error) {
	var ret []*Merchant

	agency := modules.IsAgencyAdmin(ctx)

	equalData := make(map[string]string)
	if agency != nil { // 是机构管理员的话，就需要添加机构排查
		equalData["agency_id"] = strconv.FormatUint(uint64(agency.ID), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	var userBean *account.UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*account.UserBean)
	} else {
		return 0, ret, errors.New("can't get user")
	}

	tmpDb := db.Table("merchant").Model(&Merchant{}).Where(sqlCondition)
	if (userBean.Role == string(conf.RoleUser)) && (agency == nil) { // 普通员工，需要添加账户关联
		tmpDb = tmpDb.Joins("JOIN merchant_user_associate ass ON ass.merchant_id = merchant.id AND ass.user_id = ? AND ass.deleted_at IS NULL", userBean.ID)
	}

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	if err = tmpDb.Order("updated_at desc").Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
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
