package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"time"
	"tpayment/conf"
)

var (
	db *gorm.DB
)

func InitDB() error {
	var err error
	db, err = gorm.Open(conf.DbType, os.Getenv(conf.DbConnectInfoTag))
	if err != nil {
		panic("open Db fail->" + err.Error())
	}
	db.DB().SetMaxOpenConns(256)
	db.DB().SetMaxIdleConns(8)
	db.DB().SetConnMaxLifetime(360 * time.Second)
	db.LogMode(false)

	return nil
}

func DB() *gorm.DB {
	return db
}
