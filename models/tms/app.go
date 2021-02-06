package tms

import (
	"strconv"
	"tpayment/models"
	"tpayment/modules"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type App struct {
	models.BaseModel

	AgencyId    uint64 `gorm:"column:agency_id" json:"agency_id"`
	Name        string `gorm:"column:name" json:"name"`
	PackageId   string `gorm:"column:package_id" json:"package_id"`
	Description string `gorm:"column:description" json:"description"`
}

func (App) TableName() string {
	return "tms_app"
}

// 根据device ID获取设备信息
func GetAppByID(db *models.MyDB, ctx *gin.Context, id uint64) (*App, error) {

	ret := new(App)

	err := db.Model(&App{}).Where("id=?", id).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

func QueryAppRecord(db *models.MyDB, ctx *gin.Context, offset, limit uint64, filters map[string]string) (uint64, []*App, error) {
	agency := modules.IsAgencyAdmin(ctx)

	equalData := make(map[string]string)
	if agency != nil { // 是机构管理员的话，就需要添加机构排查
		equalData["agency_id"] = strconv.FormatUint(uint64(agency.ID), 10)
	}
	sqlCondition := models.CombQueryCondition(equalData, filters)

	// conditions
	tmpDb := db.Model(&App{}).Where(sqlCondition)

	// 统计总数
	var total uint64 = 0
	err := tmpDb.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	var ret []*App
	if err = tmpDb.Order("id desc").Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		return total, ret, err
	}

	return total, ret, nil
}
