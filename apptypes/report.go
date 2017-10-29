package apptypes

import (
	"fmt"
	"time"
)

// Outcome constants //
const (
	OutcomeSuccess Outcome = iota
	OutcomeError
)

// Outcome //
type Outcome uint64

// TimingEvent //
type TimingEvent struct {
	Name                  string
	UnixNanoTimestamp     uint64
	TimestampMilliseconds float64
}

func (te TimingEvent) String() string {
	return fmt.Sprintf(
		"{Name: %s, UnixNanoTimeStamp: %d, TimestampMilliseconds: %f}",
		te.Name,
		te.UnixNanoTimestamp,
		te.TimestampMilliseconds,
	)
}

// OutcomeText //
func OutcomeText(outcome Outcome) string {
	switch outcome {
	case OutcomeSuccess:
		return "Outcome Success"
	case OutcomeError:
		return "Outcome Error"
	default:
		return "Outcome Unknown"

	}
}

// Report //
type Report struct {
	ThreadID     uint64
	RequestID    uint64
	Outcome      Outcome
	Error        string
	Code         uint64
	Response     string
	Uncompressed bool
	Timings      []TimingEvent
	Timestamp    time.Time
}
