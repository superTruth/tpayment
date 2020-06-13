package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"time"
	"tpayment/conf"
)

var (
	db *MyDB
)

type MyDB struct {
	*gorm.DB
}

func InitDB() error {
	var err error

	db = new(MyDB)

	db.DB, err = gorm.Open(conf.DbType, os.Getenv(conf.DbConnectInfoTag))
	if err != nil {
		panic("open Db fail->" + err.Error())
	}
	db.DB.DB().SetMaxOpenConns(256)
	db.DB.DB().SetMaxIdleConns(8)
	db.DB.DB().SetConnMaxLifetime(360 * time.Second)
	db.DB.LogMode(true)

	return nil
}

func DB() *MyDB {
	return db
}
