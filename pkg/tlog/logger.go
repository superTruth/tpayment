package tlog

import (
	"sync"
	"tpayment/conf"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//type Logger struct {
//	logger *zap.Logger
//	sugar  *zap.SugaredLogger
//	tag    string
//}
//
//func (this *Logger) Init(tag string) {
//	this.tag = tag
//	this.logger, _ = zap.NewProduction()
//	this.logger = this.logger.With(zap.String(conf.HeaderTagRequestId, this.tag))
//	this.sugar = this.logger.Sugar()
//}
//
//func (this *Logger) SetTag(key, value string) {
//	this.logger = this.logger.With(zap.String(key, value))
//	this.sugar = this.logger.Sugar()
//}
//
//
//func (this *Logger) Info(args ...interface{}) {
//	this.sugar.Info(args...)
//}
//
//func (this *Logger) Warn(args ...interface{}) {
//	this.sugar.Warn(args...)
//}
//
//func (this *Logger) Error(args ...interface{}) {
//	this.sugar.Error(args...)
//}
//
func SetLogger(ctx *gin.Context, logger *Logger) {
	ctx.Set(conf.ContextTagLog, logger)
}

func GetLogger(ctx *gin.Context) *Logger {
	logger, ok := ctx.Get(conf.ContextTagLog)

	if ok {
		return logger.(*Logger)
	}

	// 生成log
	log := NewLog("")

	return log
}

var once sync.Once
var logger *zap.Logger

type Logger struct {
	*zap.SugaredLogger
	Name string
}

func NewLog(tag string) *Logger {
	once.Do(func() {
		logger, _ = zap.NewProduction()
	})
	logger = logger.With(zap.String(conf.HeaderTagRequestId, tag))

	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}

func (l *Logger) Destroy() {
	//this.logger.Sync()
}
