package merchant

import (
	"strconv"
	"tpayment/models"
	"tpayment/models/account"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserMerchantAssociate struct {
	models.BaseModel

	MerchantId uint64 `json:"merchant_id" gorm:"column:merchant_id"`
	UserId     uint64 `json:"user_id" gorm:"column:user_id"`
	Role       string `json:"role" gorm:"column:role"`
}

func (UserMerchantAssociate) TableName() string {
	return "merchant_user_associate"
}

func GetUserMerchantAssociateById(db *models.MyDB, ctx *gin.Context, id uint64) (*UserMerchantAssociate, error) {
	ret := new(UserMerchantAssociate)

	err := db.Model(&UserMerchantAssociate{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetUserMerchantAssociateByMerchantIdAndUserId(db *models.MyDB, ctx *gin.Context, merchantId, userId uint64) (*UserMerchantAssociate, error) {
	ret := new(UserMerchantAssociate)

	err := db.Model(&UserMerchantAssociate{}).Where("merchant_id=? AND user_id=?", merchantId, userId).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

type AssociateMerchantUserBean struct {
	models.BaseModel

	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func QueryUsersByMerchantId(db *models.MyDB, ctx *gin.Context, merchantId, offset, limit uint64, filters map[string]string) (uint64, []*AssociateMerchantUserBean, error) {
	equalData := make(map[string]string)
	equalData["merchant_id"] = strconv.FormatUint(uint64(merchantId), 10)
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Table(account.UserBean{}.TableName()).Model(&account.UserBean{}).Where(sqlCondition).Order("id desc")
	tmpDb = tmpDb.Joins("JOIN merchant_user_associate ass ON ass.merchant_id = ? AND ass.user_id = user.id AND ass.deleted_at IS NULL", merchantId)

	// 统计总数
	var total uint64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*AssociateMerchantUserBean
	if err = tmpDb.Offset(offset).Limit(limit).Select("ass.id as id, ass.created_at as created_at, ass.updated_at as updated_at, user.name as name, user.email as email, ass.role as role").Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
