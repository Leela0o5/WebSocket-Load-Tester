package engine

import (
	"context"
	"time"

	"github.com/Leela0o5/WebSocket-Load-Tester/metrics"
	"github.com/gorilla/websocket"
)

func Worker(ctx context.Context, id int, url string, rl *RateLimiter, results chan<- metrics.Result) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		results <- metrics.Result{Err: err}
		return
	}
	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
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

		err := conn.WriteMessage(websocket.TextMessage, []byte("ping"))
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
