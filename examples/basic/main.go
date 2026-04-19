package main

import (
	"fmt"
	"time"

	"github.com/Leela0o5/LeeGo"
)

func main() {
	cfg := leego.Config{
		URL:        "ws://localhost:8080/ws",
		NumWorkers: 10,
		Duration:   5 * time.Second,
		RateLimit:  1000,
		Burst:      500,
		Message:    "ping",
	}

	fmt.Println("Starting load test...")

	stats := leego.Run(cfg)

	fmt.Printf("Test Completed!\n")
	fmt.Printf("Total Requests: %d\n", stats.TotalRequests)
	fmt.Printf("Success/Failure: %d/%d\n", stats.SuccessCount, stats.FailureCount)
	fmt.Printf("Average Latency: %s\n", stats.Average())
	fmt.Printf("P99 Latency: %s\n", stats.P99())
}
