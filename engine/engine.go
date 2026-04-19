package engine

import (
	"context"
	"sync"

	"github.com/Leela0o5/WebSocket-Load-Tester/config"
	"github.com/Leela0o5/WebSocket-Load-Tester/metrics"
)

func Run(cfg config.Config) *metrics.Stats {
	stats, done := RunAsync(cfg)
	<-done
	return stats
}

func RunAsync(cfg config.Config) (*metrics.Stats, chan struct{}) {
	stats := metrics.NewStats()
	done := make(chan struct{})

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Duration)
		defer cancel()

		rate := cfg.RateLimit
		if rate <= 0 {
			rate = cfg.NumWorkers * 1000
		}

		rl := NewRateLimiter(ctx, rate, cfg.Burst)
		results := make(chan metrics.Result, cfg.NumWorkers*10)

		collectorDone := make(chan struct{})
		go func() {
			for res := range results {
				stats.TotalRequests++
				if res.Err != nil {
					stats.FailureCount++
				} else {
					stats.SuccessCount++
					stats.Record(res.Latency)
				}
			}
			close(collectorDone)
		}()

		var wg sync.WaitGroup
		for i := 0; i < cfg.NumWorkers; i++ {
			wg.Add(1)
			id := i
			go func() {
				defer wg.Done()
				Worker(ctx, id, cfg.URL, cfg.Message, rl, results)
			}()
		}

		<-ctx.Done()
		wg.Wait()
		close(results)
		<-collectorDone
		close(done)
	}()

	return stats, done
}
