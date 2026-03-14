package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/sjingzhao/gong/pkg/httpx"
	"github.com/sjingzhao/gong/pkg/logx"
)

func Hi(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logx.FromContext(ctx).InfoContext(ctx, "http request", slog.String("name", r.PathValue("name")))
	_ = httpx.Json(w, 200, "ok", map[string]string{"trace_id": ctx.Value("trace_id").(string), "name": r.PathValue("name")})
}

// PanicTest intentionally causes a slice bounds out of range panic
func PanicTest(w http.ResponseWriter, r *http.Request) {
	var a []int
	_ = a[10] // This will panic
}

// TimeoutTest intentionally sleeps longer than the timeout middleware
func TimeoutTest(w http.ResponseWriter, r *http.Request) {
	// Our global read/write timeouts are around 10s, and maybe we add a 3s timeout middleware
	time.Sleep(5 * time.Second)
	_ = httpx.Json(w, 200, "ok", "Should have timed out!")
}
