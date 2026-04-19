package reporter

import (
	"fmt"
	"github.com/Leela0o5/LeeGo/metrics"
)

func PrintSummary(s *metrics.Stats) {
	printData(s.TotalRequests, s.FailureCount, s.Average(), s.P95())
}
func PrintReport(r *JSONReport) {
	printData(r.Total, r.Errors, r.AvgLatency, r.P95)
}

func printData(total, errors int, avg, p95 any) {
	fmt.Println("\n--- BENCHMARK SUMMARY ---")
	fmt.Printf("Total Requests: %d\n", total)
	fmt.Printf("Errors:         %d\n", errors)
	fmt.Printf("Avg Latency:    %v\n", avg)
	fmt.Printf("P95 Latency:    %v\n", p95)
	fmt.Println("-------------------------")
}
