package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewCors() gin.HandlerFunc {
	return cors.Default()
}
