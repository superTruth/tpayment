package account

import (
	"strconv"
	"tpayment/models"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type UserBean struct {
	models.BaseModel
	AgencyId uint `gorm:"column:agency_id" json:"agency_id,omitempty"`

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
func GetUserByEmail(db *models.MyDB, ctx echo.Context, email string) (*UserBean, error) {
	ret := new(UserBean)

	err := db.Model(&UserBean{}).Where("email=?", email).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetUserById(db *models.MyDB, ctx echo.Context, id uint) (*UserBean, error) {
	ret := new(UserBean)

	err := db.Model(&UserBean{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryUserRecord(db *models.MyDB, ctx echo.Context, offset, limit, agencyId uint, filters map[string]string) (uint, []*UserBean, error) {
	var ret []*UserBean

	equalData := make(map[string]string)
	if agencyId != 0 {
		equalData["agency_id"] = strconv.FormatUint(uint64(agencyId), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	tmpDB := db.Model(&UserBean{}).Where(sqlCondition)

	// 统计总数
	var total uint = 0
	err := tmpDB.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	// 查询记录
	err = tmpDB.Offset(offset).Limit(limit).Find(&ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return 0, ret, nil
		}
		return 0, nil, err
	}

	for i := 0; i < len(ret); i++ {
		ret[i].Pwd = ""
	}

	return total, ret, nil
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
func GetAppIdByAppID(db *models.MyDB, ctx echo.Context, appId string) (*AppIdBean, error) {
	ret := new(AppIdBean)

	err := db.Model(&AppIdBean{}).Where("app_id=?", appId).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetAppIdByID(db *models.MyDB, ctx echo.Context, id uint) (*AppIdBean, error) {
	ret := new(AppIdBean)

	err := db.Model(&AppIdBean{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

type TokenBean struct {
	models.BaseModel

	UserId uint   `gorm:"column:user_id" json:"user_id,omitempty"`
	AppId  uint   `gorm:"column:app_id" json:"app_id,omitempty"`
	Token  string `gorm:"column:token" json:"token,omitempty"`
}

func (TokenBean) TableName() string {
	return "user_token"
}

func GetTokenByUserId(db *models.MyDB, ctx echo.Context, userId, appId uint) (*TokenBean, error) {
	ret := new(TokenBean)

	err := db.Model(&TokenBean{}).Where("user_id=? AND app_id=?", userId, appId).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func GetTokenBeanByToken(db *models.MyDB, ctx echo.Context, token string) (*TokenBean, error) {
	ret := new(TokenBean)

	err := db.Model(&TokenBean{}).Where("token=?", token).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

type RoleBean struct {
	models.BaseModel
	Name string `gorm:"column:name" json:"name,omitempty"`
}

func (RoleBean) TableName() string {
	return "user_role"
}
