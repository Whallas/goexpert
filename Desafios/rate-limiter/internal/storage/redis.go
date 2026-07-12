package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStrategy persists counters and block flags in Redis.
// It satisfies Strategy and is the default backend for the challenge.
type RedisStrategy struct {
	client *redis.Client
}

// NewRedisStrategy builds a RedisStrategy and verifies connectivity so a
// misconfigured Redis fails fast at startup rather than on first request.
func NewRedisStrategy(ctx context.Context, addr, password string, db int) (*RedisStrategy, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connecting to redis at %q: %w", addr, err)
	}

	return &RedisStrategy{client: client}, nil
}

// Increment uses INCR (atomic) and sets the TTL only on the first hit of a
// window. Checking the returned value == 1 avoids a separate EXISTS round-trip.
func (s *RedisStrategy) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	count, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("incrementing counter %q: %w", key, err)
	}

	if count == 1 {
		if err := s.client.Expire(ctx, key, window).Err(); err != nil {
			return 0, fmt.Errorf("setting ttl on counter %q: %w", key, err)
		}
	}

	return count, nil
}

func (s *RedisStrategy) Block(ctx context.Context, key string, duration time.Duration) error {
	if err := s.client.Set(ctx, key, "1", duration).Err(); err != nil {
		return fmt.Errorf("setting block flag %q: %w", key, err)
	}
	return nil
}

func (s *RedisStrategy) IsBlocked(ctx context.Context, key string) (bool, error) {
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("checking block flag %q: %w", key, err)
	}
	return exists > 0, nil
}

// Close releases the underlying Redis connection pool.
func (s *RedisStrategy) Close() error {
	return s.client.Close()
}
