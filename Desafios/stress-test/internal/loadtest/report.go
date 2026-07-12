package loadtest

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Report aggregates the outcome of a load test run.
type Report struct {
	Elapsed      time.Duration
	Total        int
	Errors       int         // requests that never produced an HTTP response
	StatusCounts map[int]int // HTTP status code -> count
}

// Print writes a human-readable summary of the report to w.
func (r Report) Print(w io.Writer) {
	fmt.Fprintln(w, "==================== Load Test Report ====================")
	fmt.Fprintf(w, "Total time taken:          %s\n", r.Elapsed.Round(time.Millisecond))
	fmt.Fprintf(w, "Total requests performed:  %d\n", r.Total)
	fmt.Fprintf(w, "HTTP 200 responses:        %d\n", r.StatusCounts[200])

	fmt.Fprintln(w, "Status code distribution:")
	for _, code := range r.sortedCodes() {
		fmt.Fprintf(w, "  HTTP %d: %d\n", code, r.StatusCounts[code])
	}
	if r.Errors > 0 {
		fmt.Fprintf(w, "  Failed (no response): %d\n", r.Errors)
	}
	fmt.Fprintln(w, "=========================================================")
}

func (r Report) sortedCodes() []int {
	codes := make([]int, 0, len(r.StatusCounts))
	for code := range r.StatusCounts {
		codes = append(codes, code)
	}
	sort.Ints(codes)
	return codes
}
