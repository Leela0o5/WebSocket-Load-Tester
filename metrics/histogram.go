package metrics

import (
	"sort"
	"time"
)

var bucketBounds = []time.Duration{
	100 * time.Microsecond,
	200 * time.Microsecond,
	500 * time.Microsecond,
	1 * time.Millisecond,
	2 * time.Millisecond,
	5 * time.Millisecond,
	10 * time.Millisecond,
	20 * time.Millisecond,
	50 * time.Millisecond,
	100 * time.Millisecond,
	200 * time.Millisecond,
	500 * time.Millisecond,
	1 * time.Second,
	2 * time.Second,
	5 * time.Second,
	10 * time.Second,
}

type Histogram struct {
	counts  [17]uint64    
	total   uint64        
	sum     time.Duration 
	min     time.Duration 
	max     time.Duration 
}

func NewHistogram() *Histogram {
	return &Histogram{}
}

func (h *Histogram) Record(d time.Duration) {
	idx := sort.Search(len(bucketBounds), func(i int) bool {
		return bucketBounds[i] > d
	})
	h.counts[idx]++
	h.total++
	h.sum += d

	if h.total == 1 {
		h.min = d
		h.max = d
	} else {
		if d < h.min {
			h.min = d
		}
		if d > h.max {
			h.max = d
		}
	}
}

func (h *Histogram) Count() uint64 { return h.total }
func (h *Histogram) Min() time.Duration {
	if h.total == 0 {
		return 0
	}
	return h.min
}

func (h *Histogram) Max() time.Duration {
	if h.total == 0 {
		return 0
	}
	return h.max
}

func (h *Histogram) Mean() time.Duration {
	if h.total == 0 {
		return 0
	}
	return h.sum / time.Duration(h.total)
}

func (h *Histogram) Percentile(q float64) time.Duration {
	if h.total == 0 || q < 0 || q > 1 {
		return 0
	}

	target := uint64(q * float64(h.total))
	if target >= h.total {
		target = h.total - 1
	}

	var cumulative uint64
	for i, count := range h.counts {
		cumulative += count
		if cumulative > target {
			lo := lowerBound(i)
			hi := upperBound(i)
			if count == 1 {
				return lo
			}
			bucketStart := cumulative - count
			posInBucket := float64(target-bucketStart) / float64(count)
			return lo + time.Duration(posInBucket*float64(hi-lo))
		}
	}

	return h.max
}

func (h *Histogram) P50() time.Duration { return h.Percentile(0.50) }
func (h *Histogram) P95() time.Duration { return h.Percentile(0.95) }
func (h *Histogram) P99() time.Duration { return h.Percentile(0.99) }

func lowerBound(i int) time.Duration {
	if i == 0 {
		return 0
	}
	return bucketBounds[i-1]
}
func upperBound(i int) time.Duration {
	if i < len(bucketBounds) {
		return bucketBounds[i]
	}
	return bucketBounds[len(bucketBounds)-1] * 2
}
