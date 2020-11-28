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

	h, err := router.Init()
	if err != nil {
		return
	}
	if err = h.Run(":80"); err != nil {
		panic(err)
	}
}
