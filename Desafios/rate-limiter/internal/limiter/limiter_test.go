package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/Whallas/goexpert/Desafios/rate-limiter/internal/storage"
)

func TestAllowWithinLimit(t *testing.T) {
	l := New(storage.NewMemoryStrategy(), time.Minute)
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		allowed, err := l.Allow(ctx, "ip:1.1.1.1", 5)
		if err != nil {
			t.Fatalf("request %d: unexpected error: %v", i, err)
		}
		if !allowed {
			t.Fatalf("request %d within limit 5 should be allowed", i)
		}
	}
}

func TestAllowExceedsLimitThenBlocks(t *testing.T) {
	l := New(storage.NewMemoryStrategy(), time.Minute)
	ctx := context.Background()
	const limit = 3

	// First `limit` requests pass.
	for i := 1; i <= limit; i++ {
		if allowed, _ := l.Allow(ctx, "ip:2.2.2.2", limit); !allowed {
			t.Fatalf("request %d should be allowed", i)
		}
	}

	// The (limit+1)th request trips the limit and triggers the block.
	if allowed, _ := l.Allow(ctx, "ip:2.2.2.2", limit); allowed {
		t.Fatal("request over the limit should be denied")
	}

	// While blocked, even a fresh window stays denied.
	if allowed, _ := l.Allow(ctx, "ip:2.2.2.2", limit); allowed {
		t.Fatal("request during block window should be denied")
	}
}

func TestBlockExpiresAfterDuration(t *testing.T) {
	store := storage.NewMemoryStrategy()
	l := New(store, 50*time.Millisecond)
	ctx := context.Background()

	// Exceed the limit to get blocked.
	l.Allow(ctx, "ip:3.3.3.3", 1)
	if allowed, _ := l.Allow(ctx, "ip:3.3.3.3", 1); allowed {
		t.Fatal("second request should be blocked")
	}

	// After the block + window expire, requests are allowed again.
	time.Sleep(1100 * time.Millisecond)
	if allowed, _ := l.Allow(ctx, "ip:3.3.3.3", 1); !allowed {
		t.Fatal("request after block expiry should be allowed")
	}
}
