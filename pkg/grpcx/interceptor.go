package grpcx

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/sjingzhao/gong/pkg/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TraceIdInterceptor 服务端 Unary 拦截器：从 Metadata 提取 → 注入 Context
func TraceIdInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// 从 Metadata 获取客户端传来的认证信息
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	// 提取 trace_id
	var traceID string
	if ids := md.Get("x-trace-id"); len(ids) > 0 {
		traceID = ids[0]
	} else {
		traceID = uuid.New().String()
	}

	ctx = context.WithValue(ctx, "pkg", "grpc")
	ctx = context.WithValue(ctx, "full_handler", info.FullMethod)
	ctx = context.WithValue(ctx, "trace_id", traceID)

	log := defaultServerLogger.Logger
	log = log.With(slog.String("trace_id", traceID))
	log = log.With(slog.String("full_handler", info.FullMethod))
	ctx = context.WithValue(ctx, logx.CtxLoggerKey, log)

	return handler(ctx, req)
}

// RecoveryInterceptor is a unary interceptor that catches panics and returns an internal error
func RecoveryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("gRPC request panic recovered",
				"error", r,
				"method", info.FullMethod,
				"stack", string(debug.Stack()),
			)
			err = status.Errorf(codes.Internal, "Internal Server Error")
		}
	}()
	return handler(ctx, req)
}

// TimeoutInterceptor adds a timeout to the incoming gRPC request
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Only apply timeout if client didn't already set a shorter one
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		return handler(ctx, req)
	}
}
