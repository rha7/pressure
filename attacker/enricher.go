package attacker

import (
	"math"
	"sort"

	"github.com/rha7/pressure/apptypes"
)

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

func calculateResponseTimePercentils(summary *apptypes.Summary) {
	reportCount := len(summary.Reports)
	durations := make([]float64, 0, reportCount)
	for _, report := range summary.Reports {
		durations = append(durations, report.DurationMilliseconds)
	}
	sort.Float64s(durations)
	summary.Stats.ResponseTimesPercentiles.Percentile10 = durations[int(float64(reportCount)*float64(0.10))]
	summary.Stats.ResponseTimesPercentiles.Percentile25 = durations[int(float64(reportCount)*float64(0.25))]
	summary.Stats.ResponseTimesPercentiles.Percentile50 = durations[int(float64(reportCount)*float64(0.50))]
	summary.Stats.ResponseTimesPercentiles.Percentile75 = durations[int(float64(reportCount)*float64(0.75))]
	summary.Stats.ResponseTimesPercentiles.Percentile90 = durations[int(float64(reportCount)*float64(0.90))]
	summary.Stats.ResponseTimesPercentiles.Percentile95 = durations[int(float64(reportCount)*float64(0.95))]
	summary.Stats.ResponseTimesPercentiles.Percentile99 = durations[int(float64(reportCount)*float64(0.99))]
}

func calculateEventsAggregations(summary *apptypes.Summary) {
	availableEvents := make(map[apptypes.ReqEvt]bool)
	summary.Stats.EventsAggregations = make(map[apptypes.ReqEvt]*apptypes.SummaryEventsAggregationsDetail)
	for _, report := range summary.Reports {
		for _, timing := range report.Timings {
			availableEvents[timing.Name] = true
			if summary.Stats.EventsAggregations[timing.Name] == nil {
				summary.Stats.EventsAggregations[timing.Name] =
					&apptypes.SummaryEventsAggregationsDetail{
						Slowest: 0.0,
						Average: 0.0,
						Fastest: math.MaxFloat64,
					}
			}
			if summary.Stats.EventsAggregations[timing.Name].Slowest < timing.TimestampMilliseconds {
				summary.Stats.EventsAggregations[timing.Name].Slowest = timing.TimestampMilliseconds
			}
			if summary.Stats.EventsAggregations[timing.Name].Fastest > timing.TimestampMilliseconds {
				summary.Stats.EventsAggregations[timing.Name].Fastest = timing.TimestampMilliseconds
			}
			summary.Stats.EventsAggregations[timing.Name].Average += timing.TimestampMilliseconds
		}
	}
	for reqEvt := range availableEvents {
		summary.Stats.EventsAggregations[reqEvt].Average =
			summary.Stats.EventsAggregations[reqEvt].Average /
				float64(len(summary.Reports))
	}
}

func enrichSummary(summary *apptypes.Summary) {
	calculateAggregations(summary)
	calculateHTTPStatusCodeDistribution(summary)
	calculateResponseTimeDistribution(summary)
	calculateResponseTimePercentils(summary)
	calculateEventsAggregations(summary)
}
