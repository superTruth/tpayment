package tlog

import (
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"tpayment/conf"
)

type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	tag    string
}

func (this *Logger) Init(tag string) {
	this.tag = tag
	this.logger, _ = zap.NewProduction()
	this.logger = this.logger.With(zap.String("Request ID", this.tag))
	this.sugar = this.logger.Sugar()
}

func (this *Logger) Destroy() {
	//this.logger.Sync()
}

func (this *Logger) Info(args ...interface{}) {
	this.sugar.Info(args...)
}

func (this *Logger) Warn(args ...interface{}) {
	this.sugar.Warn()
}

func (this *Logger) Error(args ...interface{}) {
	this.sugar.Error(args...)
}

func SetLogger(ctx echo.Context, logger *Logger) {
	ctx.Set(conf.ContextTagLog, logger)
}

func GetLogger(ctx echo.Context) *Logger {
	logger, ok := ctx.Get(conf.ContextTagLog).(*Logger)

	if ok {
		return logger
	}

	// 生成log
	log := new(Logger)
	log.Init("")

	return log
}
