package engine

import (
	"context"
	"time"
)

type RateLimiter struct {
	tokens chan struct{}
}

func NewRateLimiter(ctx context.Context, rate, burst int) *RateLimiter {
	if burst <= 0 {
		burst = rate
	}

	rl := &RateLimiter{
		tokens: make(chan struct{}, burst),
	}

	go func() {
		const refillHz = 100
		interval := time.Second / refillHz
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		tokensPerTick := float64(rate) / float64(refillHz)
		var tokenAccumulator float64

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				tokenAccumulator += tokensPerTick
				toAdd := int(tokenAccumulator)
				tokenAccumulator -= float64(toAdd)

				for i := 0; i < toAdd; i++ {
					select {
					case rl.tokens <- struct{}{}:
					default:
						goto nextTick
					}
				}
			nextTick:
			}
		}
	}()

	return rl
}
func (rl *RateLimiter) Wait(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-rl.tokens:
		return
	}
}
