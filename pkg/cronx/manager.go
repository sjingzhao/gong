package cronx

import (
	"context"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/sjingzhao/gong/pkg/logx"
)

type Manager struct {
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	location *time.Location

	entries []*Entry
	stop    chan struct{}
	nextID  int
}

// Schedule 调度接口
type Schedule interface {
	Next(time.Time) time.Time
}

// Entry 单个定时任务
type Entry struct {
	ID       int
	Schedule Schedule
	Next     time.Time
	Prev     time.Time
	Func     Handle
}

var defaultManagerLogger = logx.NewFileLog("./logs/cron.log")

// New creates a new script manager
func New() *Manager {
	return &Manager{
		wg:       sync.WaitGroup{},
		entries:  make([]*Entry, 0),
		location: time.Local,
		stop:     make(chan struct{}),
	}
}

// Start starts all scripts concurrently
func (m *Manager) Start(ctx context.Context) error {
	m.ctx, m.cancel = context.WithCancel(ctx)
	slog.Info("Cron Manager starting", slog.Int("entries", len(m.entries)))

	now := time.Now().In(m.location)
	for _, entry := range m.entries {
		entry.Next = entry.Schedule.Next(now)
	}

	for {
		sort.Slice(m.entries, func(i, j int) bool {
			return m.entries[i].Next.Before(m.entries[j].Next)
		})

		var timer *time.Timer
		if len(m.entries) == 0 || m.entries[0].Next.IsZero() {
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			timer = time.NewTimer(m.entries[0].Next.Sub(now))
		}

		select {
		case now = <-timer.C:
			now = now.In(m.location)
			for _, e := range m.entries {
				if e.Next.After(now) || e.Next.IsZero() {
					break
				}
				m.wg.Go(func() {
					defer func() {
						if err := recover(); err != nil {
							defaultManagerLogger.Error("cron panic recovered", "error", err)
						}
					}()
					e.Func(traceIdWithContext(ctx))
				})
				e.Prev = e.Next
				e.Next = e.Schedule.Next(now)
			}

		case <-m.stop:
			timer.Stop()
			return nil
		}
	}
}

// Shutdown stops all scripts
func (m *Manager) Shutdown(ctx context.Context) error {
	close(m.stop)
	m.cancel()
	m.wg.Wait()
	return nil
}
