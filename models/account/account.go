package account

import (
	"errors"
	"strconv"
	"time"
	"tpayment/models"

	"gorm.io/gorm"
)

var UserBeanDao = &UserBean{}

type UserBean struct {
	models.BaseModel
	AgencyId uint64 `gorm:"column:agency_id" json:"agency_id,omitempty"`

	Email  string `gorm:"column:email" json:"email,omitempty"`
	Pwd    string `gorm:"column:pwd" json:"pwd,omitempty"`
	Name   string `gorm:"column:name" json:"name,omitempty"`
	Role   string `gorm:"column:role" json:"role,omitempty"`
	Active bool   `gorm:"column:active" json:"active,omitempty"`
}

func (UserBean) TableName() string {
	return "user"
}

// 通过email查询
func GetUserByEmail(email string) (*UserBean, error) {
	ret := new(UserBean)

	err := models.DB.Model(&UserBean{}).Where("email=?", email).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetUserById(id uint64) (*UserBean, error) {
	ret := new(UserBean)

	err := models.DB.Model(&UserBean{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryUserRecord(offset, limit, agencyId uint64, filters map[string]string) (uint64, []*UserBean, error) {
	var ret []*UserBean

	equalData := make(map[string]string)
	if agencyId != 0 {
		equalData["agency_id"] = strconv.FormatUint(uint64(agencyId), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	tmpDB := models.DB.Model(&UserBean{}).Where(sqlCondition)

	// 统计总数
	var total int64 = 0
	err := tmpDB.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	// 查询记录
	err = tmpDB.Order("id desc").Offset(int(offset)).Limit(int(limit)).Find(&ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return 0, ret, nil
		}
		return 0, nil, err
	}

	for i := 0; i < len(ret); i++ {
		ret[i].Pwd = ""
	}

	return uint64(total), ret, nil
}

func DeleteUser(user *UserBean) error {
	// 只是增加一下后缀，不实际删除
	data := time.Now().Unix()
	user.Email = user.Email + "-" + strconv.FormatInt(data, 10)
	return models.DB.Model(user).Select("email").Updates(user).Error
}

func (u *UserBean) UpdatePwd() error {
	tmpDB := models.DB.Model(u).Select("pwd").Updates(u)
	err := tmpDB.Error
	if err != nil {
		return err
	}
	if tmpDB.RowsAffected == 0 {
		return errors.New("no record update")
	}

	return nil
}
func (u *UserBean) GetByEmail(agencyID uint64, email string) (*UserBean, error) {
	ret := new(UserBean)

	err := models.DB.Model(ret).Where("email=? and agency_id=?", email, agencyID).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

type AppIdBean struct {
	models.BaseModel

	AppId     string `gorm:"column:app_id" json:"app_id,omitempty"`
	AppSecret string `gorm:"column:app_secret" json:"app_secret,omitempty"`
}

func (AppIdBean) TableName() string {
	return "user_app_id"
}

// 查询AppID
func GetAppIdByAppID(appId string) (*AppIdBean, error) {
	ret := new(AppIdBean)

	err := models.DB.Model(&AppIdBean{}).Where("app_id=?", appId).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetAppIdByID(id uint64) (*AppIdBean, error) {
	ret := new(AppIdBean)

	err := models.DB.Model(&AppIdBean{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

var TokenBeanDao = &TokenBean{}

type TokenBean struct {
	models.BaseModel

	UserId uint64 `gorm:"column:user_id" json:"user_id,omitempty"`
	AppId  uint64 `gorm:"column:app_id" json:"app_id,omitempty"`
	Token  string `gorm:"column:token" json:"token,omitempty"`
}

func (TokenBean) TableName() string {
	return "user_token"
}

func GetTokenByUserId(userId, appId uint64) (*TokenBean, error) {
	ret := new(TokenBean)

	err := models.DB.Model(&TokenBean{}).Where("user_id=? AND app_id=?", userId, appId).First(ret).Error

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

	err := models.DB.Model(&TokenBean{}).Where("token=?", token).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func (t *TokenBean) IsUserLogined(userId uint64) (bool, error) {
	ret := new(TokenBean)

	err := models.DB.Model(t).Where("user_id=?", userId).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return false, nil
		}
		return false, err
	}

	return true, nil
}

type RoleBean struct {
	models.BaseModel
	Name string `gorm:"column:name" json:"name,omitempty"`
}

func (RoleBean) TableName() string {
	return "user_role"
}
