package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Whallas/goexpert/Desafios/rate-limiter/internal/limiter"
	"github.com/Whallas/goexpert/Desafios/rate-limiter/internal/storage"
)

// okHandler is the protected endpoint behind the middleware.
var okHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
})

// newTestMiddleware wires the real limiter to an in-memory store so tests run
// without Redis. ipLimit=2, tokenLimit=5, token "vip" overridden to 10.
func newTestMiddleware() *RateLimiter {
	l := limiter.New(storage.NewMemoryStrategy(), time.Minute)
	return NewRateLimiter(l, 2, 5, map[string]int{"vip": 10})
}

func do(t *testing.T, h http.Handler, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestIPLimitReturns429WithExactBody(t *testing.T) {
	h := newTestMiddleware().Handler(okHandler)

	// ipLimit is 2: first two pass.
	for i := 1; i <= 2; i++ {
		if rec := do(t, h, nil); rec.Code != http.StatusOK {
			t.Fatalf("request %d: got %d, want 200", i, rec.Code)
		}
	}

	rec := do(t, h, nil)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("over-limit request: got %d, want 429", rec.Code)
	}
	body, _ := io.ReadAll(rec.Body)
	if string(body) != TooManyRequestsMessage {
		t.Fatalf("body = %q, want exact spec message", string(body))
	}
}

// TestTokenOverridesIP is the precedence test (Token > IP). The same IP that is
// limited to 2 req/s makes 5 requests carrying a token whose limit is higher,
// and all succeed — proving the token limit wins over the IP limit.
func TestTokenOverridesIP(t *testing.T) {
	h := newTestMiddleware().Handler(okHandler)
	headers := map[string]string{"API_KEY": "regular-token"} // tokenLimit = 5

	for i := 1; i <= 5; i++ {
		rec := do(t, h, headers)
		if rec.Code != http.StatusOK {
			t.Fatalf("token request %d: got %d, want 200 (token limit 5 > ip limit 2)", i, rec.Code)
		}
	}

	// 6th exceeds the token limit.
	if rec := do(t, h, headers); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("token request 6: got %d, want 429", rec.Code)
	}
}

// TestPerTokenOverride verifies a specific token can have its own higher limit.
func TestPerTokenOverride(t *testing.T) {
	h := newTestMiddleware().Handler(okHandler)
	headers := map[string]string{"API_KEY": "vip"} // overridden to 10

	for i := 1; i <= 10; i++ {
		if rec := do(t, h, headers); rec.Code != http.StatusOK {
			t.Fatalf("vip request %d: got %d, want 200", i, rec.Code)
		}
	}
	if rec := do(t, h, headers); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("vip request 11: got %d, want 429", rec.Code)
	}
}

// failingAllower simulates a broken backend to assert fail-closed behaviour.
type failingAllower struct{}

func (failingAllower) Allow(_ context.Context, _ string, _ int) (bool, error) {
	return false, context.DeadlineExceeded
}

func TestBackendErrorFailsClosed(t *testing.T) {
	m := NewRateLimiter(failingAllower{}, 2, 5, nil)
	rec := do(t, m.Handler(okHandler), nil)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("backend error: got %d, want 500", rec.Code)
	}
}
