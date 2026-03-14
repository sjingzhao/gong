package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/sjingzhao/gong/pkg/logx"
)

func ExampleScriptHandle(ctx context.Context) {
	for {
		if c := isClose(ctx); c {
			break
		}

		// do something
		logx.FromContext(ctx).InfoContext(ctx, "script process", slog.String("do", "something"))
		time.Sleep(5 * time.Second)
	}
}

func isClose(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		logx.FromContext(ctx).InfoContext(ctx, "ExampleScriptImpl.Process.isClose")
		return true
	default:
		return false
	}
}
