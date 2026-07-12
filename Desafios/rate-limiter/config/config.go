package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds every tunable, all sourced from environment variables so the
// limiter can be reconfigured without code changes (challenge requirement).
type Config struct {
	ServerPort string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// IPRateLimit is the default max requests/second per IP.
	IPRateLimit int
	// TokenRateLimit is the default max requests/second for any token without
	// an explicit override in TokenLimits.
	TokenRateLimit int
	// TokenLimits maps a token to its own requests/second limit, overriding
	// both TokenRateLimit and IPRateLimit (Token > IP precedence rule).
	TokenLimits map[string]int

	// BlockDuration is how long an offending IP/token stays blocked.
	BlockDuration time.Duration
}

// Load reads configuration from the environment, applying sane defaults so the
// app boots out-of-the-box. Returns an error only on malformed values.
func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		TokenLimits:    map[string]int{},
	}

	var err error
	if cfg.RedisDB, err = getEnvInt("REDIS_DB", 0); err != nil {
		return nil, err
	}
	if cfg.IPRateLimit, err = getEnvInt("IP_RATE_LIMIT", 10); err != nil {
		return nil, err
	}
	if cfg.TokenRateLimit, err = getEnvInt("TOKEN_RATE_LIMIT", 100); err != nil {
		return nil, err
	}

	blockSeconds, err := getEnvInt("BLOCK_DURATION_SECONDS", 300)
	if err != nil {
		return nil, err
	}
	cfg.BlockDuration = time.Duration(blockSeconds) * time.Second

	if cfg.TokenLimits, err = parseTokenLimits(getEnv("TOKEN_LIMITS", "")); err != nil {
		return nil, err
	}

	return cfg, nil
}

// parseTokenLimits parses the "token:limit,token2:limit2" format used by the
// TOKEN_LIMITS env var into a lookup map.
func parseTokenLimits(raw string) (map[string]int, error) {
	limits := map[string]int{}
	if strings.TrimSpace(raw) == "" {
		return limits, nil
	}

	for _, pair := range strings.Split(raw, ",") {
		token, value, ok := strings.Cut(strings.TrimSpace(pair), ":")
		if !ok {
			return nil, fmt.Errorf("invalid TOKEN_LIMITS entry %q: expected token:limit", pair)
		}
		limit, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return nil, fmt.Errorf("invalid limit in TOKEN_LIMITS entry %q: %w", pair, err)
		}
		limits[strings.TrimSpace(token)] = limit
	}
	return limits, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) (int, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid integer for %s=%q: %w", key, v, err)
	}
	return n, nil
}
