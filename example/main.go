package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	gong "github.com/sjingzhao/gong/core"
	"github.com/sjingzhao/gong/pb/example"
	"github.com/sjingzhao/gong/pkg/cronx"
	"github.com/sjingzhao/gong/pkg/grpcx"
	"github.com/sjingzhao/gong/pkg/httpx"
	"github.com/sjingzhao/gong/pkg/scriptx"
	"google.golang.org/grpc"
)

func main() {
	// 创建并启动应用
	app := gong.New(context.Background())

	app.WithRunnable(grpcSvr())
	app.WithRunnable(httpSvr())
	app.WithRunnable(cron())
	app.WithRunnable(script())

	if err := app.Run(); err != nil {
		slog.Error("application exited with error", "error", err)
		os.Exit(1)
	}
}

func httpSvr() *httpx.Server {
	s := httpx.NewBaseHttpServer()
	httpx.RegisterGetHandleFunc(s, "/hi/{name}", Hi)
	httpx.RegisterGetHandleFunc(s, "/panic", PanicTest)
	httpx.RegisterGetHandleFunc(s, "/timeout", TimeoutTest)

	return httpx.NewServer("0.0.0.0:8090", s)
}

func grpcSvr() *grpcx.Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcx.TraceIdInterceptor,
		grpcx.RecoveryInterceptor,
		grpcx.TimeoutInterceptor(10*time.Second),
	))

	example.RegisterUserServiceServer(s, &UserServiceImpl{})
	return grpcx.NewServer("0.0.0.0:8091", s)
}

func script() *scriptx.Manager {
	manager := scriptx.New()
	scriptx.RegisterHandle(manager, ExampleScriptHandle)
	return manager
}

func cron() *cronx.Manager {
	manager := cronx.New()
	cronx.RegisterHandle(manager, "0 * * * * *", func(ctx context.Context) { slog.Info("每分钟执行") })
	cronx.RegisterHandle(manager, "@daily", func(ctx context.Context) { slog.Info("每天执行") })
	return manager
}
