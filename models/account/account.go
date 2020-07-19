package account

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"tpayment/models"
)

type UserBean struct {
	gorm.Model
	AgencyId uint `gorm:"column:agency_id"`

	Email  string `gorm:"column:email"`
	Pwd    string `gorm:"column:pwd"`
	Name   string `gorm:"column:name"`
	Role   string `gorm:"column:role"`
	Active bool   `gorm:"column:active"`
}

func (UserBean) TableName() string {
	return "user"
}

// 通过email查询
func GetUserByEmail(email string) (*UserBean, error) {
	ret := new(UserBean)

	err := models.DB().Model(&UserBean{}).Where("email=?", email).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetUserById(id uint) (*UserBean, error) {
	ret := new(UserBean)

	err := models.DB().Model(&UserBean{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryUserRecord(offset, limit uint, filters map[string]string) (uint, []UserBean, error) {

	filterTmp := make(map[string]interface{})

	for k, v := range filters {
		filterTmp[k] = v
	}

	fmt.Println("filterTmp->", filterTmp)

	// 统计总数
	var total uint = 0
	err := models.DB().Model(&UserBean{}).Where(filterTmp).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}
	fmt.Println("total->", total)

	// 查询记录
	var ret []UserBean

	err = models.DB().Model(&UserBean{}).Where(filterTmp).Offset(offset).Limit(limit).Find(&ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return 0, ret, nil
		}
		return 0, nil, err
	}

	for _, k := range ret {
		k.Pwd = "******"
	}

	return total, ret, nil
}

func GetUsersById(ids []uint) ([]UserBean, error) {
	var ret []UserBean

	err := models.DB().Model(&UserBean{}).Where(ids).First(&ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

type AppIdBean struct {
	gorm.Model

	AppId     string `gorm:"column:app_id"`
	AppSecret string `gorm:"column:app_secret"`
}

func (AppIdBean) TableName() string {
	return "user_app_id"
}

// 查询AppID
func GetAppIdByAppID(appId string) (*AppIdBean, error) {
	ret := new(AppIdBean)

	err := models.DB().Model(&AppIdBean{}).Where("app_id=?", appId).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetAppIdByID(id uint) (*AppIdBean, error) {
	ret := new(AppIdBean)

	err := models.DB().Model(&AppIdBean{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

type TokenBean struct {
	gorm.Model

	UserId uint   `gorm:"column:user_id"`
	AppId  uint   `gorm:"column:app_id"`
	Token  string `gorm:"column:token"`
}

func (TokenBean) TableName() string {
	return "user_token"
}

func GetTokenByUserId(userId, appId uint) (*TokenBean, error) {
	ret := new(TokenBean)

	err := models.DB().Model(&TokenBean{}).Where("user_id=? AND app_id=?", userId, appId).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetTokenBeanByToken(token string) (*TokenBean, error) {
	ret := new(TokenBean)

	err := models.DB().Model(&TokenBean{}).Where("token=?", token).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

type RoleBean struct {
	gorm.Model

	Name string `gorm:"column:name"`
}

func (RoleBean) TableName() string {
	return "user_role"
}
