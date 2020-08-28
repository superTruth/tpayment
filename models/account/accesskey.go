package account

import (
	"errors"
	"strconv"
	"tpayment/conf"
	"tpayment/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserAccessKeyBean struct {
	models.BaseModel

	UserId uint   `gorm:"column:user_id" json:"user_id,omitempty"`
	Key    string `gorm:"column:key" json:"key,omitempty"`
	Secret string `gorm:"column:secret" json:"secret,omitempty"`
}

func (UserAccessKeyBean) TableName() string {
	return "user_access_key"
}

func GetUserAccessKeyFromID(db *models.MyDB, ctx *gin.Context, id uint) (*UserAccessKeyBean, error) {
	ret := new(UserAccessKeyBean)

	err := db.Model(&UserAccessKeyBean{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetUserAccessKeyFromKey(db *models.MyDB, ctx *gin.Context, key string) (*UserAccessKeyBean, error) {
	ret := new(UserAccessKeyBean)

	err := db.Model(&UserAccessKeyBean{}).Where("key=?", key).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryAccessKeysRecord(db *models.MyDB, ctx *gin.Context, offset, limit uint, filters map[string]string) (uint, []*UserAccessKeyBean, error) {
	var ret []*UserAccessKeyBean

	equalData := make(map[string]string)

	var userBean *UserBean
	userBeanTmp, ok := ctx.Get(conf.ContextTagUser)
	if ok {
		userBean = userBeanTmp.(*UserBean)
	} else {
		return 0, ret, errors.New("can't get user")
	}

	equalData["user_id"] = strconv.FormatUint(uint64(userBean.ID), 10)
	sqlCondition := models.CombQueryCondition(equalData, filters)

	tmpDB := db.Model(&UserAccessKeyBean{}).Where(sqlCondition)

	// 统计总数
	var total uint = 0
	err := tmpDB.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	// 查询记录
	err = tmpDB.Order("updated_at desc").Offset(offset).Limit(limit).Find(&ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return 0, ret, nil
		}
		return 0, nil, err
	}

	for i := 0; i < len(ret); i++ {
		ret[i].Secret = ""
	}

	return total, ret, nil
}
