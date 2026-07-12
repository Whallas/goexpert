package middleware

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
)

// TooManyRequestsMessage is the exact body required by the challenge spec.
const TooManyRequestsMessage = "you have reached the maximum number of requests or actions allowed within a certain time frame"

// apiKeyHeader carries the access token (challenge: "API_KEY: <TOKEN>").
const apiKeyHeader = "API_KEY"

// Allower is the slice of limiter behaviour the middleware depends on. Keeping
// it an interface decouples the HTTP layer from the limiter implementation.
type Allower interface {
	Allow(ctx context.Context, key string, limit int) (bool, error)
}

// RateLimiter is the HTTP middleware. It resolves the identity (token or IP)
// and its limit, then delegates the decision to the injected Allower.
type RateLimiter struct {
	limiter        Allower
	ipLimit        int
	tokenLimit     int
	tokenOverrides map[string]int
}

func NewRateLimiter(l Allower, ipLimit, tokenLimit int, tokenOverrides map[string]int) *RateLimiter {
	if tokenOverrides == nil {
		tokenOverrides = map[string]int{}
	}
	return &RateLimiter{
		limiter:        l,
		ipLimit:        ipLimit,
		tokenLimit:     tokenLimit,
		tokenOverrides: tokenOverrides,
	}
}

// Handler wraps next with rate limiting.
func (m *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, limit := m.resolve(r)

		allowed, err := m.limiter.Allow(r.Context(), key, limit)
		if err != nil {
			// Fail closed: a broken backend must not silently disable limiting.
			log.Printf("rate limiter error for key %q: %v", key, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(TooManyRequestsMessage))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// resolve picks the rate-limit identity and its limit. Token settings override
// IP settings (Token > IP precedence). The token's own limit is its override
// if present, otherwise the global token limit.
func (m *RateLimiter) resolve(r *http.Request) (key string, limit int) {
	if token := r.Header.Get(apiKeyHeader); token != "" {
		if override, ok := m.tokenOverrides[token]; ok {
			return "token:" + token, override
		}
		return "token:" + token, m.tokenLimit
	}
	return "ip:" + clientIP(r), m.ipLimit
}

// clientIP extracts the caller IP, preferring the proxy-supplied header so the
// limiter keys on the real client rather than a shared proxy address.
func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		// X-Forwarded-For may be a comma list; the first entry is the client.
		first := strings.TrimSpace(strings.Split(fwd, ",")[0])
		if host, _, err := net.SplitHostPort(first); err == nil {
			return host
		}
		return first
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
