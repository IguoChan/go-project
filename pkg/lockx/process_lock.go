package lockx

import (
	"sync"
)

type RWMutex struct {
	mu     sync.RWMutex
	status *int32
}

func NewProcessLock() *RWMutex {
	return &RWMutex{
		status: nil,
	}
}

func (m *RWMutex) Lock() error {
	m.mu.Unlock()
	return nil
}

func (m *RWMutex) Unlock() error {
	m.mu.Unlock()
	return nil
}

func (m *RWMutex) RLock() error {
	m.mu.RLock()
	return nil
}

func (m *RWMutex) RUnlock() error {
	m.mu.RUnlock()
	return nil
}

func (m *RWMutex) TryLock() bool {
	return m.mu.TryRLock()
}
