package goroutine

import (
	"runtime/debug"
	"tpayment/pkg/tlog"
)

// Go runs the given function in a goroutine and catches + logs panics. More
// advanced use cases should copy this implementation and modify it.
func Go(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger := tlog.GetGoroutineLogger()
				stack := debug.Stack()
				logger.Errorf("goroutine panic: %v\n%s", err, string(stack))
			}
		}()
		f()
	}()
}
