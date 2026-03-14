package cronx

import "context"

type Handle func(ctx context.Context)

func RegisterHandle(m *Manager, spec string, handler Handle) {
	schedule, err := Parse(spec)
	if err != nil {
		defaultManagerLogger.Logger.Error("parse spec failed", "error", err)
		return
	}
	m.nextID++
	entry := &Entry{
		ID:       m.nextID,
		Schedule: schedule,
		Func:     handler,
	}
	m.entries = append(m.entries, entry)
}
