package scriptx

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/sjingzhao/gong/pkg/logx"
)

func traceIdWithContext(ctx context.Context) context.Context {
	var traceID = uuid.New().String()

	ctx = context.WithValue(ctx, "pkg", "script")
	ctx = context.WithValue(ctx, "trace_id", traceID)

	log := defaultManagerLogger.Logger
	log = log.With(slog.String("trace_id", traceID))
	ctx = context.WithValue(ctx, logx.CtxLoggerKey, log)

	return ctx
}
