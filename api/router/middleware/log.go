package middleware

import (
	"bytes"
	"net/http/httputil"
	"tpayment/conf"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Logger(ctx *gin.Context) {
	// 生成request ID
	requestId := uuid.New().String()

	// 生成log
	//logger := new(tlog.Logger)
	//logger.Init(requestId)

	logger := tlog.NewLog(requestId)

	tlog.SetLogger(ctx, logger)
	defer logger.Destroy()
	tlog.SetGoroutineLogger(logger)
	defer tlog.FreeGoroutineLogger()

	ctx.Request.Header[conf.HeaderTagRequestId] = []string{requestId}
	ctx.Writer.Header()[conf.HeaderTagRequestId] = []string{requestId}

	content, _ := httputil.DumpRequest(ctx.Request, true)
	logger.Info("request->", string(content))

	// 执行，并且拷贝response对象
	writer := &responseBodyWriter{bodyCache: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
	ctx.Writer = writer

	ctx.Next()

	responseDump := writer.dump()
	logger.Info("Response->", string(responseDump))
}

type responseBodyWriter struct {
	gin.ResponseWriter
	bodyCache *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.bodyCache.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w responseBodyWriter) WriteString(s string) (int, error) {
	w.bodyCache.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
func (w responseBodyWriter) dump() []byte {
	return w.bodyCache.Bytes()
}
