package printers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
)

// Text //
func Text(logger *logrus.Logger, spec apptypes.TestSpec, reports []apptypes.Report) error {
	var (
		totalTimeMin float64
		totalTimeAvg float64
		totalTimeMax float64
		totalTimes   []float64
	)
	fmt.Printf("Timing Report\n")
	fmt.Println()
	fmt.Printf("Test Spec:\n")
	fmt.Printf("%s", indentTextBlock(yamlPrettyPrint(spec), "  "))
	fmt.Printf("Report Count: %d\n", len(reports))
	fmt.Println()
	for idxReport, report := range reports {
		totalTimeMilliseconds := report.Timings[len(report.Timings)-1].TimestampMilliseconds
		fmt.Printf("Report #%d\n", idxReport+1)
		fmt.Printf("\tCode: %d\n", report.Code)
		fmt.Printf("\tCodeText: %s\n", http.StatusText(int(report.Code)))
		fmt.Printf("\tError: %s\n", report.Error)
		fmt.Printf("\tOutcome: %s\n", apptypes.OutcomeText(report.Outcome))
		fmt.Printf("\tCompressed: %t\n", report.Compressed)
		fmt.Printf("\tTimestamp: %s\n", report.Timestamp)
		fmt.Printf("\tTotalTimeMilliseconds: %f\n", totalTimeMilliseconds)
		fmt.Printf("\tTimings:\n")
		for idxEvent, event := range report.Timings {
			fmt.Printf("\t\tSequence: %d\n", idxEvent+1)
			fmt.Printf("\t\t\tName: %s\n", event.Name)
			fmt.Printf("\t\t\tTimestamp: %v\n", time.Unix(0, int64(event.UnixNanoTimestamp)))
			fmt.Printf("\t\t\tTimestampMilliseconds: %f\n", event.TimestampMilliseconds)
		}
		totalTimes = append(totalTimes, totalTimeMilliseconds)
	}
	totalTimeMin = totalTimes[0]
	totalTimeAvg = 0
	totalTimeMax = totalTimes[0]
	for _, t := range totalTimes {
		if t < totalTimeMin {
			totalTimeMin = t
		}
		totalTimeAvg += t
		if t > totalTimeMax {
			totalTimeMax = t
		}
	}
	totalTimeAvg = totalTimeAvg / float64(len(totalTimes))
	fmt.Printf("Statistics\n")
	fmt.Printf("\tMinimum: %fms\n", totalTimeMin)
	fmt.Printf("\tAverage: %fms\n", totalTimeAvg)
	fmt.Printf("\tMaximum: %fms\n", totalTimeMax)
	fmt.Println()
	return nil
}
