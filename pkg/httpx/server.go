package httpx

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"

	"github.com/sjingzhao/gong/pkg/logx"
)

var defaultServerLogger = logx.NewFileLog("./logs/http.log")

// Server is an http server
type Server struct {
	addr string
	s    *http.Server
}

// ServerOption is a function type that configures a Server
type ServerOption func(server *Server)

// NewServer creates a new Server
func NewServer(addr string, s *http.Server, opts ...ServerOption) *Server {
	svr := &Server{addr: addr, s: s}
	for _, opt := range opts {
		opt(svr)
	}
	return svr
}

// Start starts the server and blocks until shutdown
func (server *Server) Start(ctx context.Context) error {
	l, err := net.Listen("tcp", server.addr)
	if err != nil {
		slog.Error("http listen error", "error", err)
		return err
	}
	slog.Info("HTTP server starting", "addr", l.Addr().String())



	if err = server.s.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Shutdown stops the server
func (server *Server) Shutdown(ctx context.Context) error {
	if server.s == nil {
		return nil
	}

	if err := server.s.Shutdown(ctx); err != nil {
		return err
	}

	return defaultServerLogger.Close()
}
