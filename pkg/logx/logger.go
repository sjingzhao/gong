package logx

import (
	"io"
	"log/slog"
	"os"
)

// Logger 日志记录器，封装了 slog.Logger
type Logger struct {
	w            io.WriteCloser // 日志写入器（文件或标准输出）
	*slog.Logger                // 底层日志记录器
}

// New 创建一个新的日志记录器
// file: 日志文件路径，如果失败则使用标准输出
func New(w io.WriteCloser) *Logger {
	opt := &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}
	handler := slog.NewJSONHandler(w, opt)
	log := &Logger{w: w, Logger: slog.New(handler)}
	return log
}

func NewFileLog(file string) *Logger {
	stdout, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", "error", err)
		stdout = os.Stdout
	}

	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	return &Logger{w: stdout, Logger: slog.New(slog.NewJSONHandler(stdout, opts))}
}

// With 添加键值对到日志记录器，返回新的实例
func (l *Logger) With(args ...any) *slog.Logger {
	return l.Logger.With(args...)
}

// Close 关闭日志记录器，释放资源
func (l *Logger) Close() error {
	return safeClose(l.w)
}

// safeClose 安全地关闭一个 io.Writer
// 如果是标准输出或标准错误，不会真正关闭
func safeClose(w io.Writer) error {
	// 如果它是空的，直接返回
	if w == nil {
		return nil
	}

	// 核心：如果是标准输出或标准错误，绝对不要 Close！直接放行
	if w == os.Stdout || w == os.Stderr {
		return nil
	}

	// 如果它实现了 io.Closer 接口，才调用 Close
	if closer, ok := w.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
