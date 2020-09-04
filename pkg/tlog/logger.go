package tlog

import (
	"sync"
	"tpayment/conf"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
	tmpLogger := logger.With(zap.String(conf.HeaderTagRequestId, tag))

	return &Logger{
		SugaredLogger: tmpLogger.Sugar(),
	}
}

func (l *Logger) Destroy() {
	//this.logger.Sync()
}
