package goroutine

import (
	"runtime/debug"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

// Go runs the given function in a goroutine and catches + logs panics. More
// advanced use cases should copy this implementation and modify it.
func Go(f func(), ctx *gin.Context) {
	go func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger := tlog.GetGoroutineLogger()
				stack := debug.Stack()
				logger.Error("goroutine panic: %v\n%s", err, string(stack))
			}
		}()
		f()
	}(ctx)
}
