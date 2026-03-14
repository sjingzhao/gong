package gong

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Runnable represents a lifecycle component that can be started and stopped
type Runnable interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// RunnableManager manages runnable components
type RunnableManager struct {
	runnables []Runnable
}

// newRunnableManager creates a new runnable manager
func newRunnableManager() *RunnableManager {
	return &RunnableManager{runnables: make([]Runnable, 0)}
}

// add adds a runnable to the manager
func (m *RunnableManager) add(runnable Runnable) {
	m.runnables = append(m.runnables, runnable)
}

// start starts all runnables concurrently using errgroup
func (m *RunnableManager) start(ctx context.Context) (*errgroup.Group, context.Context) {
	eg, egCtx := errgroup.WithContext(ctx)

	for _, runnable := range m.runnables {
		r := runnable
		eg.Go(func() error {
			return r.Start(egCtx)
		})
	}
	
	// 返回 errgroup 引用让调用层可以 Wait
	return eg, egCtx
}

// shutdown stops all runnables concurrently
func (m *RunnableManager) shutdown(ctx context.Context) error {
	var shutdownGroup errgroup.Group

	for _, runnable := range m.runnables {
		r := runnable
		shutdownGroup.Go(func() error {
			return r.Shutdown(ctx)
		})
	}

	return shutdownGroup.Wait()
}
