package loadtest

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Config holds the parameters of a load test run.
type Config struct {
	URL         string
	Requests    int
	Concurrency int
	Timeout     time.Duration
}

// Validate checks that the configuration is usable.
func (c Config) Validate() error {
	if c.URL == "" {
		return errors.New("--url is required")
	}
	if u, err := url.ParseRequestURI(c.URL); err != nil || u.Scheme == "" || u.Host == "" {
		return errors.New("--url must be a valid absolute URL")
	}
	if c.Requests <= 0 {
		return errors.New("--requests must be greater than zero")
	}
	if c.Concurrency <= 0 {
		return errors.New("--concurrency must be greater than zero")
	}
	if c.Concurrency > c.Requests {
		return errors.New("--concurrency cannot exceed --requests")
	}
	return nil
}

// result carries the outcome of a single request.
type result struct {
	status int // 0 means the request never completed (transport error)
}

// Run executes the load test described by cfg and returns its report.
func Run(cfg Config) Report {
	// Allow the connection pool to keep one idle connection per worker so the
	// keep-alive connections are actually reused across requests.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = cfg.Concurrency
	transport.MaxIdleConnsPerHost = cfg.Concurrency

	client := &http.Client{Timeout: cfg.Timeout, Transport: transport}

	jobs := make(chan struct{})
	results := make(chan result)

	var wg sync.WaitGroup
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go worker(client, cfg.URL, jobs, results, &wg)
	}

	start := time.Now()

	// Feed exactly cfg.Requests jobs to the worker pool.
	go func() {
		for i := 0; i < cfg.Requests; i++ {
			jobs <- struct{}{}
		}
		close(jobs)
	}()

	// Close results once every worker has drained the jobs channel.
	go func() {
		wg.Wait()
		close(results)
	}()

	report := Report{StatusCounts: make(map[int]int)}
	for r := range results {
		report.Total++
		if r.status == 0 {
			report.Errors++
			continue
		}
		report.StatusCounts[r.status]++
	}
	report.Elapsed = time.Since(start)

	return report
}

func worker(client *http.Client, target string, jobs <-chan struct{}, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	for range jobs {
		results <- doRequest(client, target)
	}
}

func doRequest(client *http.Client, target string) result {
	resp, err := client.Get(target)
	if err != nil {
		return result{status: 0}
	}
	defer resp.Body.Close()
	// Drain the body so the connection can be reused (keep-alive). Without this
	// every request opens a fresh TCP/TLS connection, which is slow and makes
	// throttling hosts (e.g. google.com) appear to hang under high request counts.
	io.Copy(io.Discard, resp.Body)
	return result{status: resp.StatusCode}
}
