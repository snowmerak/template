package httpserver

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	return &Server{
		router: gin.Default(),
	}
}

const defaultCorsMaxAge = 12 * time.Hour

func (s *Server) EnableCORS(origins []string, methods []string, headers []string) {
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     methods,
		AllowHeaders:     headers,
		AllowCredentials: true,
		MaxAge:           defaultCorsMaxAge,
	}))
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) GET(relativePath string, handlers ...gin.HandlerFunc) {
	s.router.GET(relativePath, handlers...)
}

func (s *Server) POST(relativePath string, handlers ...gin.HandlerFunc) {
	s.router.POST(relativePath, handlers...)
}

func (s *Server) PUT(relativePath string, handlers ...gin.HandlerFunc) {
	s.router.PUT(relativePath, handlers...)
}

func (s *Server) DELETE(relativePath string, handlers ...gin.HandlerFunc) {
	s.router.DELETE(relativePath, handlers...)
}

func (s *Server) PATCH(relativePath string, handlers ...gin.HandlerFunc) {
	s.router.PATCH(relativePath, handlers...)
}

func (s *Server) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) {
	s.router.OPTIONS(relativePath, handlers...)
}

func (s *Server) HEAD(relativePath string, handlers ...gin.HandlerFunc) {
	s.router.HEAD(relativePath, handlers...)
}

func (s *Server) Static(relativePath, root string) {
	s.router.Static(relativePath, root)
}

func (s *Server) StaticFile(relativePath, filepath string) {
	s.router.StaticFile(relativePath, filepath)
}

func (s *Server) Use(handlers ...gin.HandlerFunc) {
	s.router.Use(handlers...)
}
