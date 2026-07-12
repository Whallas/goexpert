package storage

import (
	"context"
	"sync"
	"time"
)

// MemoryStrategy is an in-process Strategy implementation. It exists to prove
// the Strategy pattern: the limiter runs against it with zero changes, and the
// unit tests use it to avoid a Redis dependency.
//
// Not for production (state is per-process and lost on restart).
type MemoryStrategy struct {
	mu       sync.Mutex
	counters map[string]counter
	blocks   map[string]time.Time
	now      func() time.Time // injectable clock for deterministic tests
}

type counter struct {
	value     int64
	expiresAt time.Time
}

func NewMemoryStrategy() *MemoryStrategy {
	return &MemoryStrategy{
		counters: make(map[string]counter),
		blocks:   make(map[string]time.Time),
		now:      time.Now,
	}
}

func (s *MemoryStrategy) Increment(_ context.Context, key string, window time.Duration) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	c, ok := s.counters[key]
	if !ok || now.After(c.expiresAt) {
		c = counter{value: 0, expiresAt: now.Add(window)}
	}

	c.value++
	s.counters[key] = c
	return c.value, nil
}

func (s *MemoryStrategy) Block(_ context.Context, key string, duration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.blocks[key] = s.now().Add(duration)
	return nil
}

func (s *MemoryStrategy) IsBlocked(_ context.Context, key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	expiresAt, ok := s.blocks[key]
	if !ok {
		return false, nil
	}
	if s.now().After(expiresAt) {
		delete(s.blocks, key)
		return false, nil
	}
	return true, nil
}
