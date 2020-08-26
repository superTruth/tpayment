package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewCors() gin.HandlerFunc {
	config := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders: []string{"X-Access-Token", "token", "DNT", "X-Mx-ReqToken", "Keep-Alive", "User-Agent",
			"X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}

	return cors.New(config)
}

func Cors(ctx *gin.Context) {
	method := ctx.Request.Method
	fmt.Println("Truth Cors")

	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "X-Access-Token,token,DNT,X-Mx-ReqToken,Keep-Alive,"+
		"User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization")
	ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	ctx.Header("Access-Control-Expose-Headers", "X-Access-Token,token,DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization")
	ctx.Header("Access-Control-Allow-Credentials", "true")

	//放行所有OPTIONS方法
	if method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
	}
	// 处理请求
	ctx.Next()
}
