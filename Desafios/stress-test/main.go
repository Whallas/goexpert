package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Whallas/goexpert/Desafios/stress-test/internal/loadtest"
)

func main() {
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		flag.Usage()
		os.Exit(1)
	}

	report := loadtest.Run(cfg)
	report.Print(os.Stdout)
}

func parseFlags() (loadtest.Config, error) {
	var cfg loadtest.Config

	flag.StringVar(&cfg.URL, "url", "", "URL of the service under test")
	flag.IntVar(&cfg.Requests, "requests", 0, "total number of requests to perform")
	flag.IntVar(&cfg.Concurrency, "concurrency", 1, "number of simultaneous calls")
	flag.DurationVar(&cfg.Timeout, "timeout", 10*time.Second, "per-request timeout")
	flag.Parse()

	return cfg, cfg.Validate()
}
