package models

import (
	"context"
	"time"
	"tpayment/conf"
	"tpayment/pkg/tlog"

	"gorm.io/driver/mysql"

	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDB() {
	var err error

	DB = new(gorm.DB)

	openDb, err := gorm.Open(mysql.Open(conf.GetConfigData().DBAccount), &gorm.Config{
		Logger: &dbLogger{},
	})
	DB = openDb

	if err != nil {
		panic("open Db fail->" + err.Error())
	}
	sqlDB, err := DB.DB()
	if err != nil {
		panic("open Db fail2->" + err.Error())
	}
	sqlDB.SetMaxOpenConns(256)
	sqlDB.SetMaxIdleConns(8)
	sqlDB.SetConnMaxLifetime(360 * time.Second)
}

type dbLogger struct {
}

func (d *dbLogger) LogMode(logger.LogLevel) logger.Interface {
	return d
}

func (d dbLogger) Info(ctx context.Context, template string, args ...interface{}) {
	tlog.GetGoroutineLogger().Infof(template, args...)
}

func (d dbLogger) Warn(ctx context.Context, template string, args ...interface{}) {
	tlog.GetGoroutineLogger().Warnf(template, args...)
}

func (d dbLogger) Error(ctx context.Context, template string, args ...interface{}) {
	tlog.GetGoroutineLogger().Errorf(template, args...)
}

const warningLimitTime = 3 * time.Second

func (d dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	log := tlog.GetGoroutineLogger()
	if time.Since(begin) > warningLimitTime { // sql语句处理时间太长，需要报警，可能产生慢查询
		log.Errorf("sql processing time too long")
	}
	sql, rows := fc()
	log.Info(sql, ", rows:", rows)
}
