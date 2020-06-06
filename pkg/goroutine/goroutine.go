package goroutine

import (
	"github.com/labstack/echo"
	"runtime/debug"
	"tpayment/pkg/tlog"
)

// Go runs the given function in a goroutine and catches + logs panics. More
// advanced use cases should copy this implementation and modify it.
func Go(f func(), ctx echo.Context) {
	go func(ctx echo.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger := tlog.GetLogger(ctx)
				stack := debug.Stack()
				logger.Error("goroutine panic: %v\n%s", err, stack)
			}
		}()
		f()
	}(ctx)
}
