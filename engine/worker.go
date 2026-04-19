package engine

import (
	"context"
	"time"

	"github.com/Leela0o5/LeeGo/metrics"
	"github.com/gorilla/websocket"
)

func Worker(ctx context.Context, id int, url string, message string, rl *RateLimiter, results chan<- metrics.Result) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		results <- metrics.Result{Err: err}
		return
	}
	defer conn.Close()

	doneCh := make(chan struct{})
	defer close(doneCh)

	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		case <-doneCh:
			return
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		rl.Wait(ctx)
		select {
		case <-ctx.Done():
			return
		default:
		}

		start := time.Now()

		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			results <- metrics.Result{
				Err:       err,
				Success:   false,
				Timestamp: time.Now(),
			}
			return
		}

		_, _, err = conn.ReadMessage()
		if err != nil {
			results <- metrics.Result{
				Err:       err,
				Success:   false,
				Timestamp: time.Now(),
			}
			return
		}

		results <- metrics.Result{
			Latency:   time.Since(start),
			Success:   true,
			Timestamp: start,
		}
	}
}
