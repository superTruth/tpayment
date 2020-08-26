package middleware

import (
	"runtime/debug"
	"tpayment/conf"
	"tpayment/modules"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

// NewRecovery creates Recovery middleware
func NewRecovery(c *gin.Context) {
	defer newRecovery(c)
	c.Next()
}

func newRecovery(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	err := recover()
	if err == nil {
		return
	}

	debug.PrintStack()

	logger.Error("exception->", string(debug.Stack()))

	modules.BaseError(ctx, conf.PanicError)
	ctx.Abort()
}
