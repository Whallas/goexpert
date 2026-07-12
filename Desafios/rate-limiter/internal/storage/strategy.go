package storage

import (
	"context"
	"time"
)

// Strategy abstracts the persistence backend used by the rate limiter.
//
// The challenge mandates Redis, but the limiter depends only on this
// interface. Swapping persistence (e.g. in-memory, Memcached, DynamoDB)
// means writing a new Strategy implementation — no limiter changes.
type Strategy interface {
	// Increment adds 1 to the counter stored at key and returns the new value.
	// The first increment in a fresh window must set the key TTL to window so
	// the counter resets automatically (fixed-window algorithm).
	Increment(ctx context.Context, key string, window time.Duration) (int64, error)

	// Block marks key as blocked for the given duration. New lookups via
	// IsBlocked must report true until the duration elapses.
	Block(ctx context.Context, key string, duration time.Duration) error

	// IsBlocked reports whether key is currently within an active block window.
	IsBlocked(ctx context.Context, key string) (bool, error)
}
