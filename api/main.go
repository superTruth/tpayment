package main

import (
	"tpayment/api/router"
	"tpayment/conf"
	"tpayment/internal/basekey"
	"tpayment/models"
)

func main() {
	conf.InitConfigData()

	models.InitDB()

	basekey.Init() // 初始化基础秘钥

	//Test2()

	h, err := router.Init()
	if err != nil {
		return
	}
	if err = h.Run(":80"); err != nil {
		panic(err)
	}
}

func Test2() {
	bean := &Test{
		BaseModel: models.BaseModel{
			ID: 0,
		},
		Name: "123",
	}
	_ = models.DB().Model(bean).Create(bean).Error
}

type Test struct {
	models.BaseModel
	Name string `gorm:"column:name"`
}

func (Test) TableName() string {
	return "test"
}
