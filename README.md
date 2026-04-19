# LeeGo: High-Performance WebSocket Load Tester

[![Go Reference](https://pkg.go.dev/badge/github.com/Leela0o5/LeeGo.svg)](https://pkg.go.dev/github.com/Leela0o5/LeeGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/Leela0o5/LeeGo)](https://goreportcard.com/report/github.com/Leela0o5/LeeGo)

LeeGo is a fast WebSocket load testing tool built in Go. It helps you measure how your servers handle concurrent WebSocket connections and tracks latency without tanking your system's memory.

## Features

- **Real Concurrency:** Uses goroutines to keep thousands of connections open and active at the same time.
- **Efficient Latency Tracking:** Uses a custom histogram bucket architecture to track P50, P95, and P99 latencies accurately without massive Garbage Collection pauses.
- **Token-Bucket Rate Limiting:** Set exact requests-per-second limits so you can find the breaking point of your servers smoothly.
- **CLI or Library:** Run it from the terminal for a nice TUI, or just import `leego` directly into your own Go codebase.

## Installation

### As a CLI Tool

Install the standalone executable onto your machine:
```bash
go install github.com/Leela0o5/LeeGo/cmd/leego@latest
```

### As a Go Library

Add it to your project dependencies:
```bash
go get github.com/Leela0o5/LeeGo
```

## Usage

### Library Usage
If you want to automate stress tests from inside your backend or CI, just drop `leego` into your Go code.

```go
package main

import (
    "fmt"
    "time"

    "github.com/Leela0o5/LeeGo" // imported as leego
)

func main() {
    cfg := leego.Config{
        URL:        "ws://localhost:8080/ws",
        NumWorkers: 50,
        Duration:   10 * time.Second,
        RateLimit:  2000,
        Message:    "ping",
    }

    stats := leego.Run(cfg)

    fmt.Printf("Total Requests: %d | Approximated P99: %s\n", stats.TotalRequests, stats.P99())
}
```

### Command-Line Usage

Run dynamic load tests against your servers straight from the terminal. 

```bash
# Basic run parameters
leego run --url ws://localhost:8080/ws --connections 100 --duration 10s --rate 5000

# Run using a configuration file and save a report
leego run -c config.yaml -o results.json
```

**Available Flags:**
- `--url` (string): WebSocket connection URL to target.
- `--connections` (int): Number of concurrent workers (simulated clients).
- `--duration` (duration): Total test duration (e.g., `10s`, `1m`).
- `--rate` (int): Rate limit per second (total requests across all connections).
- `--burst` (int): Burst size limit.
- `--message` (string): Plaintext or JSON payload message to send.
- `-c, --config` (string): Path to a YAML configuration file.
- `-o, --output` (string): Path to save a structured JSON report.

#### Configuration File Example
```yaml
url: "ws://localhost:8080/ws"
connections: 50
duration: "30s"
rate: 2000
burst: 1000
message: "{ \"type\": \"ping\" }"
```

## Tests

Execute the validation suites:
```bash
go test ./...
```
