package metrics

import (
	"slices"
	"time"
)

type Stats struct {
	TotalRequests int
	SuccessCount  int
	FailureCount  int
	Latencies     []time.Duration
	Hist          *Histogram
}

func NewStats() *Stats {
	return &Stats{
		Latencies: make([]time.Duration, 0),
		Hist:      NewHistogram(),
	}
}

func (s *Stats) Record(d time.Duration) {
	s.Latencies = append(s.Latencies, d)
	s.Hist.Record(d)
}

func (s *Stats) Average() time.Duration {
	if s.Hist != nil && s.Hist.Count() > 0 {
		return s.Hist.Mean()
	}
	return s.sliceAverage()
}

func (s *Stats) Min() time.Duration {
	if s.Hist != nil && s.Hist.Count() > 0 {
		return s.Hist.Min()
	}
	return s.sliceMin()
}

func (s *Stats) Max() time.Duration {
	if s.Hist != nil && s.Hist.Count() > 0 {
		return s.Hist.Max()
	}
	return s.sliceMax()
}

func (s *Stats) Median() time.Duration { return s.Percentile(0.50) }
func (s *Stats) P95() time.Duration    { return s.Percentile(0.95) }
func (s *Stats) P99() time.Duration    { return s.Percentile(0.99) }

func (s *Stats) Percentile(p float64) time.Duration {
	if s.Hist != nil && s.Hist.Count() > 0 {
		return s.Hist.Percentile(p)
	}
	return s.slicePercentile(p)
}

// --- Slice-based fallbacks (original implementation) ---

func (s *Stats) sliceAverage() time.Duration {
	if len(s.Latencies) == 0 {
		return 0
	}
	var total time.Duration
	for _, l := range s.Latencies {
		total += l
	}
	return total / time.Duration(len(s.Latencies))
}

func (s *Stats) sliceMin() time.Duration {
	if len(s.Latencies) == 0 {
		return 0
	}
	s.sortLatencies()
	return s.Latencies[0]
}

func (s *Stats) sliceMax() time.Duration {
	if len(s.Latencies) == 0 {
		return 0
	}
	s.sortLatencies()
	return s.Latencies[len(s.Latencies)-1]
}

func (s *Stats) slicePercentile(p float64) time.Duration {
	if len(s.Latencies) == 0 {
		return 0
	}
	s.sortLatencies()
	idx := int(float64(len(s.Latencies)) * p)
	if idx >= len(s.Latencies) {
		idx = len(s.Latencies) - 1
	}
	return s.Latencies[idx]
}

func (s *Stats) sortLatencies() {
	slices.Sort(s.Latencies)
}
