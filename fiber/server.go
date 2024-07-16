package fiber

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pires/go-proxyproto"
	"net"
)

type Server struct {
	app *fiber.App
}

func New() *Server {
	return &Server{app: fiber.New()}
}

func (s *Server) ListenAndServe(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	proxyListener := &proxyproto.Listener{Listener: lis}

	if err := s.app.Listener(proxyListener); err != nil {
		return fmt.Errorf("failed to set listener: %w", err)
	}

	return nil
}

func (s *Server) ListenAndServeTLS(addr, certFile, keyFile string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	proxyListener := &proxyproto.Listener{Listener: lis}

	if err := s.app.Server().ServeTLS(proxyListener, certFile, keyFile); err != nil {
		return fmt.Errorf("failed to set listener: %w", err)
	}

	return nil
}
