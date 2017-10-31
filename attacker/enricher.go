package attacker

import "github.com/rha7/pressure/apptypes"

// TODO: Add Response time histogram (Ten slices between slowest and fastest, with count of each)
// TODO: Latency distribution (response time percentils, 10, 25, 50, 75, 90, 95, and 99)

func calculateHTTPStatusCodeDistribution(summary *apptypes.Summary) {
	summary.Stats.HTTPStatusCodeDistribution = make(map[uint64]uint64)
	for _, report := range summary.Reports {
		summary.Stats.HTTPStatusCodeDistribution[report.Code]++
	}
}

func calculateAggregations(summary *apptypes.Summary) {
	summary.Stats = apptypes.SummaryStats{
		Aggregations: apptypes.SummaryStatsAggregations{},
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
	}
	summary.Stats.Aggregations.Average = float64(summary.Stats.Aggregations.Total / float64(summary.Stats.Aggregations.Count))
	summary.Stats.Aggregations.RequestsPerSec = float64(summary.Stats.Aggregations.Count) / (summary.Stats.Aggregations.Total / float64(1000.0))
	summary.Stats.Aggregations.ErrorsPercent = float64(100*summary.Stats.Aggregations.Errors) / float64(summary.Stats.Aggregations.Count)
}

func calculateResponseTimeDistribution(summary *apptypes.Summary) {
	rtds := make([]uint64, 11, 11)
	delta := summary.Stats.Aggregations.Slowest
	delta -= summary.Stats.Aggregations.Fastest
	delta /= 10.0
	for _, report := range summary.Reports {
		slot := int((report.DurationMilliseconds - summary.Stats.Aggregations.Fastest) / delta)
		rtds[slot]++
	}
	summary.Stats.ResponseTimesDistribution =
		make([]apptypes.SummaryStatsResponseTimesDistributionSlot, 11, 11)
	for slot, count := range rtds {
		summary.Stats.ResponseTimesDistribution[slot] =
			apptypes.SummaryStatsResponseTimesDistributionSlot{
				From:  summary.Stats.Aggregations.Fastest + (float64(slot) * delta),
				UpTo:  summary.Stats.Aggregations.Fastest + (float64(slot) * delta) + delta,
				Count: count,
			}
	}
}

func enrichSummary(summary *apptypes.Summary) {
	calculateAggregations(summary)
	calculateHTTPStatusCodeDistribution(summary)
	calculateResponseTimeDistribution(summary)
}
