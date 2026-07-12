package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Whallas/goexpert/Desafios/rate-limiter/config"
	"github.com/Whallas/goexpert/Desafios/rate-limiter/internal/limiter"
	"github.com/Whallas/goexpert/Desafios/rate-limiter/internal/middleware"
	"github.com/Whallas/goexpert/Desafios/rate-limiter/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	// .env is optional: in Docker, env comes from compose. Ignore "not found".
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	ctx := context.Background()

	// Strategy selection happens here, in one place. Swapping persistence is a
	// one-line change (e.g. storage.NewMemoryStrategy()).
	store, err := storage.NewRedisStrategy(ctx, cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("initializing redis strategy: %v", err)
	}
	defer store.Close()

	rateLimiter := limiter.New(store, cfg.BlockDuration)
	rateMiddleware := middleware.NewRateLimiter(
		rateLimiter,
		cfg.IPRateLimit,
		cfg.TokenRateLimit,
		cfg.TokenLimits,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("request accepted\n"))
	})

	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           rateMiddleware.Handler(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("rate limiter listening on :%s (ip=%d req/s, token=%d req/s, block=%s)",
		cfg.ServerPort, cfg.IPRateLimit, cfg.TokenRateLimit, cfg.BlockDuration)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(fmt.Errorf("server failed: %w", err))
	}
}
