package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/snowmerak/template/httpserver/handler"
	"github.com/snowmerak/template/httpserver/middleware"
)

func RegisterRoutes(engine *gin.Engine) {
	engine.GET("/", handler.Index)
}

func RegisterMiddlewares(engine *gin.Engine) {
	engine.Use(middleware.EnableCORS([]string{"*"}, []string{"GET", "POST", "PUT", "DELETE"}, []string{"Origin", "Content-Type"}))
}
