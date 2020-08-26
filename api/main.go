package main

import (
	"tpayment/api/router"
	"tpayment/conf"
	"tpayment/internal/encryption"
	"tpayment/models"
)

func main() {

	conf.InitConfigData()

	models.InitDB()

	encryption.Init() // 初始化基础秘钥

	h, err := router.Init()
	if err != nil {
		return
	}
	if err = h.Run(":80"); err != nil {
		panic(err)
	}
}
