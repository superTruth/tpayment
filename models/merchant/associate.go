package merchant

import (
	"github.com/jinzhu/gorm"
	"tpayment/models"
)

type UserMerchantAssociate struct {
	gorm.Model

	MerchantId uint   `json:"merchant_id" gorm:"column:merchant_id"`
	UserId     uint   `json:"user_id" gorm:"column:user_id"`
	Role       string `json:"role" gorm:"column:role"`
}

func (UserMerchantAssociate) TableName() string {
	return "merchant_user_associate"
}

func GetUserMerchantAssociateById(id uint) (*UserMerchantAssociate, error) {
	ret := new(UserMerchantAssociate)

	err := models.DB().Model(&UserMerchantAssociate{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetUserMerchantAssociateByUserId(id uint) ([]UserMerchantAssociate, error) {
	var ret []UserMerchantAssociate

	err := models.DB().Model(&UserMerchantAssociate{}).Where("user_id=?", id).Find(&ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetUserMerchantAssociateByMerchantId(id uint) ([]UserMerchantAssociate, error) {
	var ret []UserMerchantAssociate

	err := models.DB().Model(&UserMerchantAssociate{}).Where("merchant_id=?", id).Find(&ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryMerchantsByUserRecord(userId, offset, limit uint, filters map[string]string) (uint, []Merchant, error) {
	filterTmp := make(map[string]interface{})

	for k, v := range filters {
		filterTmp[k] = v
	}

	// conditions
	tmpDb := models.DB().Table("merchant").Where(filterTmp)
	tmpDb = tmpDb.Joins("JOIN merchant_user_associate ass ON ass.merchant_id = merchant.id AND ass.user_id = ? AND ass.deleted_at IS NULL", userId)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []Merchant
	if err = tmpDb.Offset(offset).Limit(limit).Select("user.id as id, ass.created_at as created_at, ass.updated_at as updated_at, user.name as name, user.email as email, ass.role as role").Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}

type AssociateMerchantUserBean struct {
	gorm.Model

	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func QueryUsersByMerchantId(merchantId, offset, limit uint, filters map[string]string) (uint, []AssociateMerchantUserBean, error) {
	filterTmp := make(map[string]interface{})

	for k, v := range filters {
		filterTmp[k] = v
	}

	// conditions
	tmpDb := models.DB().Table("user").Where(filterTmp)
	tmpDb = tmpDb.Joins("JOIN merchant_user_associate ass ON ass.merchant_id = ? AND ass.user_id = user.id AND ass.deleted_at IS NULL", merchantId)

	// 统计总数
	var total uint = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []AssociateMerchantUserBean
	if err = tmpDb.Offset(offset).Limit(limit).Select("user.id as id, ass.created_at as created_at, ass.updated_at as updated_at, user.name as name, user.email as email, ass.role as role").Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
