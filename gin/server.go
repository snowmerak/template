package gin

import (
	"github.com/gin-gonic/gin"

	"github.com/snowmerak/template/gin/handler"
	"github.com/snowmerak/template/gin/middleware"
)

func RegisterRoutes(engine *gin.Engine) {
	engine.GET("/", handler.Index)
}

func RegisterMiddlewares(engine *gin.Engine) {
	engine.Use(middleware.EnableCORS([]string{"*"}, []string{"GET", "POST", "PUT", "DELETE"}, []string{"Origin", "Content-Type"}))
}
