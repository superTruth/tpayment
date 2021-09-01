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

var Dao = &Merchant{}

type Merchant struct {
	models.BaseModel

	AgencyId uint64 `json:"agency_id" gorm:"column:agency_id"`
	Name     string `json:"name,omitempty" gorm:"column:name"`
	Tel      string `json:"tel,omitempty"  gorm:"column:tel"`
	Addr     string `json:"addr,omitempty" gorm:"column:addr"`
	Email    string `json:"email,omitempty" gorm:"column:email"`
}

func (Merchant) TableName() string {
	return "merchant"
}

func GetMerchantById(db *models.MyDB, ctx *gin.Context, id uint64) (*Merchant, error) {
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

func (m *Merchant) Get(id uint64) (*Merchant, error) {
	ret := new(Merchant)
	err := models.DB().Model(m).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

// 获取机构下面的商户
func QueryMerchantInAgency(db *models.MyDB, ctx *gin.Context, agencyId, offset, limit uint64, filters map[string]string) (uint64, []*Merchant, error) {
	equalData := make(map[string]string)
	if agencyId != 0 {
		equalData["agency_id"] = strconv.FormatUint(uint64(agencyId), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Model(&Merchant{}).Where(sqlCondition)

	// 统计总数
	var total uint64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*Merchant
	if err = tmpDb.Order("id desc").Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

// 获取账号相关的商户列表
func QueryMerchantInUser(db *models.MyDB, ctx *gin.Context, offset, limit uint64, filters map[string]string) (uint64, []*Merchant, error) {
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
	var total uint64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	if err = tmpDb.Order("id desc").Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

func QueryMerchantByDeviceID(ctx *gin.Context, deviceSn string, offset, limit uint64) (uint64, []*Merchant, error) {
	var ret []*Merchant

	agency := modules.IsAgencyAdmin(ctx)

	equalData := make(map[string]string)
	if agency != nil { // 是机构管理员的话，就需要添加机构排查
		equalData["agency_id"] = strconv.FormatUint(agency.ID, 10)
	}

	// conditions
	var userBean *account.UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*account.UserBean)
	} else {
		return 0, ret, errors.New("can't get user")
	}

	var tmpDb *gorm.DB
	if agency != nil { // 如果是agency还需要增加一个判断
		tmpDb = models.DB().Model(&Merchant{}).
			Joins("join merchant_device md join tms_device d on d.device_sn like ? and d.agency_id=? and d.deleted_at is null and "+
				"d.id = md.device_id and md.deleted_at is null and md.merchant_id = merchant.id", deviceSn+"%", agency.ID)
		tmpDb = tmpDb.Where("agency_id = ?", agency.ID)
	} else if userBean.Role == string(conf.RoleAdmin) { // 超级管理员
		tmpDb = models.DB().Model(&Merchant{}).
			Joins("join merchant_device md join tms_device d on d.device_sn like ? and d.deleted_at is null and "+
				"d.id = md.device_id and md.deleted_at is null and md.merchant_id = merchant.id", deviceSn+"%")
	} else { // 普通员工
		tmpDb = models.DB().Model(&Merchant{}).
			Joins("join merchant_device md join tms_device d on d.device_sn like ? and d.deleted_at is null and "+
				"d.id = md.device_id and md.deleted_at is null and md.merchant_id = merchant.id"+
				" JOIN merchant_user_associate ass ON ass.merchant_id = merchant.id AND ass.user_id = ? AND ass.deleted_at IS NULL",
				deviceSn+"%", userBean.ID)
	}

	// 统计总数
	var total uint64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	if err = tmpDb.Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

func GetMerchantsById(ids []uint64) ([]Merchant, error) {
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
