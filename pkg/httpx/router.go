package httpx

import (
	"context"
	"net"
	"net/http"
	"time"
)

// HttpServerOption is a function that configures an http.Server
type HttpServerOption func(server *http.Server)

// NewBaseHttpServer Create a basic HTTP server
func NewBaseHttpServer(opts ...HttpServerOption) *http.Server {
	server := &http.Server{
		Handler:           http.NewServeMux(),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		BaseContext: func(l net.Listener) context.Context {
			return context.Background()
		},
	}
	for _, opt := range opts {
		opt(server)
	}
	return server
}

// RegisterHandleFunc registers a handler function for a given HTTP method and pattern
func RegisterGetHandleFunc(svr *http.Server, pattern string, f http.HandlerFunc, middlewares ...MiddlewareFunc) {
	RegisterHandleFunc(svr, http.MethodGet+" "+pattern, f, middlewares...)
}

// RegisterHeadHandleFunc registers a handler function for HEAD requests
func RegisterHeadHandleFunc(svr *http.Server, pattern string, f http.HandlerFunc, middlewares ...MiddlewareFunc) {
	RegisterHandleFunc(svr, http.MethodHead+" "+pattern, f, middlewares...)
}

// RegisterPostHandleFunc registers a handler function for POST requests
func RegisterPostHandleFunc(svr *http.Server, pattern string, f http.HandlerFunc, middlewares ...MiddlewareFunc) {
	RegisterHandleFunc(svr, http.MethodPost+" "+pattern, f, middlewares...)
}

// RegisterPutHandleFunc registers a handler function for PUT requests
func RegisterPutHandleFunc(svr *http.Server, pattern string, f http.HandlerFunc, middlewares ...MiddlewareFunc) {
	RegisterHandleFunc(svr, http.MethodPut+" "+pattern, f, middlewares...)
}

// RegisterPatchHandleFunc registers a handler function for PATCH requests
func RegisterPatchHandleFunc(svr *http.Server, pattern string, f http.HandlerFunc, middlewares ...MiddlewareFunc) {
	RegisterHandleFunc(svr, http.MethodPatch+" "+pattern, f, middlewares...)
}

// RegisterDeleteHandleFunc registers a handler function for DELETE requests
func RegisterDeleteHandleFunc(svr *http.Server, pattern string, f http.HandlerFunc, middlewares ...MiddlewareFunc) {
	RegisterHandleFunc(svr, http.MethodDelete+" "+pattern, f, middlewares...)
}

// RegisterConnectHandleFunc registers a handler function for CONNECT requests
func RegisterConnectHandleFunc(svr *http.Server, pattern string, f http.HandlerFunc, middlewares ...MiddlewareFunc) {
	RegisterHandleFunc(svr, http.MethodConnect+" "+pattern, f, middlewares...)
}

// RegisterOptionsHandleFunc registers a handler function for OPTIONS requests
func RegisterOptionsHandleFunc(svr *http.Server, pattern string, f http.HandlerFunc, middlewares ...MiddlewareFunc) {
	RegisterHandleFunc(svr, http.MethodOptions+" "+pattern, f, middlewares...)
}

// RegisterTraceHandleFunc registers a handler function for TRACE requests
func RegisterTraceHandleFunc(svr *http.Server, pattern string, f http.HandlerFunc, middlewares ...MiddlewareFunc) {
	RegisterHandleFunc(svr, http.MethodTrace+" "+pattern, f, middlewares...)
}

// RegisterHandleFunc registers a handler function for a given HTTP method and pattern
func RegisterHandleFunc(svr *http.Server, pattern string, f http.HandlerFunc, middlewares ...MiddlewareFunc) {
	svr.Handler.(*http.ServeMux).HandleFunc(pattern, useMiddlewareFunc(f, middlewares...))
}
