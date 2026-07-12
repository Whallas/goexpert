package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/Whallas/goexpert/Desafios/rate-limiter/internal/storage"
)

// window is the fixed time bucket for the requests-per-second limit.
const window = time.Second

// Limiter holds the rate-limiting business logic. It is intentionally free of
// any HTTP/middleware concerns so it can be tested and reused in isolation.
type Limiter struct {
	strategy      storage.Strategy
	blockDuration time.Duration
}

func New(strategy storage.Strategy, blockDuration time.Duration) *Limiter {
	return &Limiter{
		strategy:      strategy,
		blockDuration: blockDuration,
	}
}

// Allow reports whether a request identified by key (an IP or token) is within
// its limit (requests per second). On exceeding the limit the key is blocked
// for blockDuration and every subsequent request returns false until it lapses.
func (l *Limiter) Allow(ctx context.Context, key string, limit int) (bool, error) {
	blocked, err := l.strategy.IsBlocked(ctx, blockKey(key))
	if err != nil {
		return false, fmt.Errorf("checking block state for %q: %w", key, err)
	}
	if blocked {
		return false, nil
	}

	count, err := l.strategy.Increment(ctx, counterKey(key), window)
	if err != nil {
		return false, fmt.Errorf("counting requests for %q: %w", key, err)
	}

	if count > int64(limit) {
		if err := l.strategy.Block(ctx, blockKey(key), l.blockDuration); err != nil {
			return false, fmt.Errorf("blocking %q: %w", key, err)
		}
		return false, nil
	}

	return true, nil
}

// Namespaced keys keep counters and block flags from colliding in the store.
func counterKey(key string) string { return "ratelimit:count:" + key }
func blockKey(key string) string   { return "ratelimit:block:" + key }
