package metrics

func Collector(results <-chan Result) Stats {
	s := NewStats()

	for res := range results {
		s.TotalRequests++

		if res.Err != nil {
			s.FailureCount++
			continue
		}

		s.SuccessCount++
		s.Record(res.Latency)
	}

	return *s
}
