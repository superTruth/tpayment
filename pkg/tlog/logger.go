package tlog

import (
	"fmt"
	"sync"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var once sync.Once
var logger *zap.Logger
var ErrChan chan string

type Logger struct {
	*zap.SugaredLogger
	Name string

	Tags map[string]string
}

func NewLog(tags map[string]string) *Logger {
	once.Do(func() {
		config := zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json",
			EncoderConfig:    myEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}

		logger, _ = config.Build()
	})

	var fields []zap.Field
	for k, v := range tags {
		fields = append(fields, zap.String(k, v))
	}

	tmpLogger := logger.With(fields...)

	return &Logger{
		SugaredLogger: tmpLogger.Sugar(),
		Tags:          tags,
	}
}

func myEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	errMsg := fmt.Sprintf(template, args...)

	requestID, ok := l.Tags["request-id"]
	if ok {
		errMsg = "request-id:" + requestID + ":" + errMsg
	}

	if ErrChan != nil && cap(ErrChan) > len(ErrChan) {
		ErrChan <- errMsg
	}

	l.SugaredLogger.Error(errMsg)
}

func (l *Logger) Error(args ...interface{}) {
	errMsg := fmt.Sprint(args...)

	requestID, ok := l.Tags["request-id"]
	if ok {
		errMsg = "request-id:" + requestID + ":" + errMsg
	}

	if ErrChan != nil && cap(ErrChan) > len(ErrChan) {
		ErrChan <- errMsg
	}

	l.SugaredLogger.Error(errMsg)
}
