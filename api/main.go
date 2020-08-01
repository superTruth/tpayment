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
	if err = h.Start(":80"); err != nil {
		panic(err)
	}
}
