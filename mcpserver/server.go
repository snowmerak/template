package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/server"
)

type Option struct {
	Name    string
	Version string
}

type Server struct {
	s *server.MCPServer
}

func NewServer(ctx context.Context, opt Option) *Server {
	s := server.NewMCPServer(opt.Name, opt.Version, server.WithLogging(), server.WithResourceCapabilities(true, true))
	return &Server{s: s}
}

func (s *Server) Start() error {
	return server.ServeStdio(s.s)
}
