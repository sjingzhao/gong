package httpx

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/sjingzhao/gong/pkg/logx"
)

// MiddlewareFunc is a middleware function
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// useMiddlewareFunc is a helper function that applies middlewares to a http.HandlerFunc
func useMiddlewareFunc(f http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {
	if len(middlewares) > 0 {
		// 从后往前包装中间件，确保执行时按正序调用
		// 例如：middlewares = [A, B, C]，最终执行顺序为 A -> B -> C
		for i := len(middlewares) - 1; i >= 0; i-- {
			f = middlewares[i](f)
		}
	}

	// 默认必定挂载的核心安全中间件
	f = recoveryMiddleware(f)
	f = timeoutMiddleware(3 * time.Second)(f)
	f = traceIdMiddleware(f)

	return f
}

// traceIdMiddleware is a middleware that extracts trace_id from request header
func traceIdMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 提取 trace_id
		var traceID string
		if id := r.Header.Get("x-trace-id"); len(id) > 0 {
			traceID = id
		} else {
			traceID = uuid.New().String()
		}

		_ = r.ParseForm()

		ctx := r.Context()
		ctx = context.WithValue(ctx, "pkg", "http")
		ctx = context.WithValue(ctx, "full_handler", r.URL.Path)
		ctx = context.WithValue(ctx, "trace_id", traceID)

		log := defaultServerLogger.Logger
		log = log.With(slog.String("trace_id", traceID))
		log = log.With(slog.String("full_handler", r.URL.Path))
		ctx = context.WithValue(ctx, logx.CtxLoggerKey, log)

		r = r.WithContext(ctx)
		next(w, r)
	}
}

// recoveryMiddleware catches any panics in the request processing and returns a 500 error
func recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logx.FromContext(r.Context()).InfoContext(
					r.Context(),
					"recoveryMiddleware PANIC",
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.Any("err", err),
					slog.String("stack", string(debug.Stack())),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// timeoutMiddleware adds a timeout to the request context
func timeoutMiddleware(timeout time.Duration) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			// 使用 channel 通知处理完成
			done := make(chan bool, 1)

			go func() {
				next(w, r)
				done <- true
			}()

			select {
			case <-done:
				// 正常处理完成
				return
			case <-ctx.Done():
				// 超时
				w.WriteHeader(http.StatusGatewayTimeout)
				_, _ = w.Write([]byte("Gateway Timeout"))
			}
		}
	}
}

// CORSMiddleware is a middleware that handles CORS
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Trace-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	}
}
