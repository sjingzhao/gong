package gong

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Application a gong application
type Application struct {
	ctx    context.Context
	cancel context.CancelFunc
	rm     *RunnableManager
}

// Option is an option
type Option func(*Application)

// New creates a new application
func New(ctx context.Context, opts ...Option) *Application {
	ctx, cancel := context.WithCancel(ctx)
	app := &Application{ctx: ctx, cancel: cancel, rm: newRunnableManager()}
	slog.Info("initializing application")

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// WithRunnable adds a runnable component to the application
func (app *Application) WithRunnable(r Runnable) *Application {
	app.rm.add(r)
	return app
}

// Run starts the application and blocks until shutdown
func (app *Application) Run() error {
	slog.Info("starting application")

	// 启动所有的组件，并接管它们的生命周期
	eg, ctx := app.rm.start(app.ctx)

	// 开一个单独的协程来监听操作系统信号
	eg.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			// 其他组件发生错误崩溃了，就不必等信号了
			_ = app.Shutdown()
			return nil
		case s := <-quit:
			slog.Info("received shutdown signal", "signal", s)
			_ = app.Shutdown()
			return nil
		}
	})

	// 阻塞在此：如果内部发生了任何致命错误（例如端口占用且没被捕捉），立刻返回
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// Shutdown stops the application
func (app *Application) Shutdown() error {
	slog.Info("stopping application")
	app.cancel()

	// Timeout for shutdown to prevent hanging
	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelTimeout()

	// 调用组件管理器的并发下线
	if err := app.rm.shutdown(timeoutCtx); err != nil {
		slog.Error("application shutdown failed", "error", err)
		return err
	}
	slog.Info("application stopped")
	return nil
}
