package main

import (
	"tpayment/api/router"
	"tpayment/conf"
	"tpayment/models"
)

func main() {

	models.InitDB()

	conf.InitConfigData()

	h, err := router.Init()
	if err != nil {
		return
	}
	h.Start(":80")
}
