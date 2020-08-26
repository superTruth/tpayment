package middleware

import (
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
