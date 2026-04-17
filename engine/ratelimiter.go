package engine

import (
	"context"
	"time"
)

type RateLimiter struct {
	tokens chan struct{}
}
func NewRateLimiter(ctx context.Context, rate int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, rate),
	}

	go func() {
		interval := time.Second / time.Duration(rate)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				select {
				case rl.tokens <- struct{}{}:
				default:
				}
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