package main

import (
	"tpayment/api/router"
	"tpayment/models"
)

func main() {

	models.InitDB()

	h, err := router.Init()
	if err != nil {
		return
	}
	h.Start(":80")
}
