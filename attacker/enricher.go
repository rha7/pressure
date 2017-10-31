package attacker

import "github.com/rha7/pressure/apptypes"

// TODO: Add Response time histogram (Ten slices between slowest and fastest, with count of each)
// TODO: Latency distribution (response time percentils, 10, 25, 50, 75, 90, 95, and 99)

func enrichSummary(summary *apptypes.Summary) {
	summary.Stats = apptypes.SummaryStats{
		Aggregations:               apptypes.SummaryStatsAggregations{},
		HTTPStatusCodeDistribution: make(map[uint64]uint64),
	}

	responseTimeOfFirst := float64(summary.Reports[0].DurationMilliseconds)
	summary.Stats.Aggregations.Count = uint64(len(summary.Reports))
	summary.Stats.Aggregations.Slowest = responseTimeOfFirst
	summary.Stats.Aggregations.Fastest = responseTimeOfFirst
	for _, report := range summary.Reports {
		reportResponseTime := float64(report.DurationMilliseconds)
		if reportResponseTime > summary.Stats.Aggregations.Slowest {
			summary.Stats.Aggregations.Slowest = reportResponseTime
		}
		if reportResponseTime < summary.Stats.Aggregations.Fastest {
			summary.Stats.Aggregations.Fastest = reportResponseTime
		}
		if report.Outcome == apptypes.OutcomeError && reportResponseTime > 0.0 {
			summary.Stats.Aggregations.Errors++
		}
		summary.Stats.Aggregations.Total += reportResponseTime
		summary.Stats.HTTPStatusCodeDistribution[report.Code]++
	}
	summary.Stats.Aggregations.Average = float64(summary.Stats.Aggregations.Total / float64(summary.Stats.Aggregations.Count))
	summary.Stats.Aggregations.RequestsPerSec = float64(summary.Stats.Aggregations.Count) / (summary.Stats.Aggregations.Total / float64(1000.0))
	summary.Stats.Aggregations.ErrorsPercent = float64(100*summary.Stats.Aggregations.Errors) / float64(summary.Stats.Aggregations.Count)
}
