package engine

import (
	"context"
	"sync"
	"github.com/Leela0o5/WebSocket-Load-Tester/config"
	"github.com/Leela0o5/WebSocket-Load-Tester/metrics"
)

func Run(cfg config.Config) *metrics.Stats {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Duration)
	defer cancel()

	rate := cfg.RateLimit
	if rate <= 0 {
		rate = cfg.NumWorkers * 1000
	}
	rl := NewRateLimiter(ctx, rate)

	stats := &metrics.Stats{}
	results := make(chan metrics.Result, cfg.NumWorkers*10)

	done := make(chan bool)
	go func() {
		for res := range results {
			stats.TotalRequests++
			if res.Err != nil {
				stats.FailureCount++
			} else {
				stats.SuccessCount++
				stats.Latencies = append(stats.Latencies, res.Latency)
			}
		}
		done <- true
	}()

	var wg sync.WaitGroup
	for i := 0; i < cfg.NumWorkers; i++ {
		wg.Add(1)
		workerID := i
		go func() {
			defer wg.Done()
			Worker(ctx, workerID, cfg.URL, rl, results)
		}()
	}

	<-ctx.Done()
	wg.Wait()
	close(results)
	<-done

	return stats
}
