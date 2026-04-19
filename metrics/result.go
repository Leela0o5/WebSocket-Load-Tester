package metrics

import (
	"time"
)

type Result struct {
	Latency   time.Duration
	Success   bool
	Timestamp time.Time
	Err       error
}
