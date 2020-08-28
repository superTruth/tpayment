package tlog

import (
	"strconv"
	"tpayment/pkg/gls"
)

var loggerMap = make(map[uint64]*Logger)

func SetGoroutineLogger(logger *Logger) {
	loggerMap[gls.GetGoroutineID()] = logger
}

func GetGoroutineLogger() *Logger {
	id := gls.GetGoroutineID()
	logger, ok := loggerMap[id]
	if !ok {
		logger = NewLog(strconv.FormatUint(id, 10))
		return logger
	}

	return logger
}

func FreeGoroutineLogger() {
	id := gls.GetGoroutineID()
	delete(loggerMap, id)
}
