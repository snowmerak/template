package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

const defaultCorsMaxAge = 12 * time.Hour

func EnableCORS(origins []string, methods []string, headers []string) func(ctx *gin.Context) {
	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     methods,
		AllowHeaders:     headers,
		AllowCredentials: true,
		MaxAge:           defaultCorsMaxAge,
	})
}
