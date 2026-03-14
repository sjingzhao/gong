package logx

import (
	"context"
	"log/slog"
	"os"
)

const CtxLoggerKey string = "logger" // 日志记录器的上下文键

// WithContext 将日志记录器添加到上下文中
// ctx: 父上下文
// key: 组件名称，会自动作为字段添加到日志中
// log: 日志记录器实例
func WithContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, CtxLoggerKey, log)
}

// FromContext 从上下文中获取日志记录器
// ctx: 上下文
// key: 组件名称（如果未找到，会返回一个默认的 logger）
// 返回值：日志记录器，如果上下文中不存在则返回默认 logger
func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(CtxLoggerKey).(*slog.Logger); ok {
		return log
	}
	// 返回一个默认 logger，避免 panic
	return New(os.Stdout).Logger
}
