package grpcx

import (
	"context"
	"log/slog"
	"net"

	"github.com/sjingzhao/gong/pkg/logx"
	"google.golang.org/grpc"
)

var defaultServerLogger = logx.NewFileLog("./logs/grpc.log")

// Server is a gRPC server
type Server struct {
	addr string
	s    *grpc.Server
}

// ServerOption is a function type that configures a Server
type ServerOption func(server *Server)

// NewServer creates a new Server
func NewServer(addr string, s *grpc.Server, opts ...ServerOption) *Server {
	svr := &Server{addr: addr, s: s}
	for _, opt := range opts {
		opt(svr)
	}
	return svr
}

// Start starts the server and blocks until shutdown or error
func (server *Server) Start(ctx context.Context) error {
	l, err := net.Listen("tcp", server.addr)
	if err != nil {
		slog.Error("grpc listen error", "error", err)
		return err
	}
	slog.Info("GRPC server starting", "addr", l.Addr().String())



	return server.s.Serve(l)
}

// Shutdown stops the server
func (server *Server) Shutdown(ctx context.Context) error {
	if server.s == nil {
		return nil
	}

	// Wait for connections to close safely. The application provides an overarching timeout in ctx.
	done := make(chan struct{})
	go func() {
		server.s.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		server.s.Stop() // Force stop if timeout occurs
		return ctx.Err()
	case <-done:
	}

	return defaultServerLogger.Close()
}
