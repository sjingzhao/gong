package scriptx

import "context"

type Handle func(ctx context.Context)

func RegisterHandle(m *Manager, handler Handle) {
	if m.handle == nil {
		m.handle = make([]Handle, 0)
	}
	m.handle = append(m.handle, handler)
}
