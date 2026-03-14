package scriptx

import (
	"context"
	"log/slog"
	"sync"

	"github.com/sjingzhao/gong/pkg/logx"
)

type Manager struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	handle []Handle
}

var defaultManagerLogger = logx.NewFileLog("./logs/script.log")

// New creates a new script manager
func New() *Manager {
	return &Manager{wg: sync.WaitGroup{}, handle: make([]Handle, 0)}
}

// Start starts all scripts concurrently
func (m *Manager) Start(ctx context.Context) error {
	m.ctx, m.cancel = context.WithCancel(ctx)
	slog.Info("Script Manager starting", slog.Int("handler", len(m.handle)))
	for _, handle := range m.handle {
		m.wg.Go(func() {
			defer func() {
				if err := recover(); err != nil {
					defaultManagerLogger.Error("script panic recovered", "error", err)
				}
			}()
			handle(traceIdWithContext(m.ctx))
		})
	}
	return nil
}

// Shutdown stops all scripts
func (m *Manager) Shutdown(ctx context.Context) error {
	m.cancel()
	m.wg.Wait()
	return nil
}
