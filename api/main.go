package main

import (
	"tpayment/api/router"
	"tpayment/conf"
	"tpayment/models"
)

func main() {

	conf.InitConfigData()

	models.InitDB()

	h, err := router.Init()
	if err != nil {
		return
	}
	h.Start(":80")
}
