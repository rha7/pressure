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
	Name                  ReqEvt
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
	ThreadID             uint64
	RequestID            uint64
	Outcome              Outcome
	Error                string
	Code                 uint64
	Response             string
	Compressed           bool
	Timings              []TimingEvent
	Timestamp            time.Time
	DurationMilliseconds float64
}

// SummaryStatsAggregations //
type SummaryStatsAggregations struct {
	Count          uint64
	Errors         uint64
	ErrorsPercent  float64
	Total          float64
	Slowest        float64
	Fastest        float64
	Average        float64
	RequestsPerSec float64
}

// SummaryStatsResponseTimesDistributionSlot //
type SummaryStatsResponseTimesDistributionSlot struct {
	From  float64
	UpTo  float64
	Count uint64
}

// Latency distribution:
//   10% in 0.3437 secs
//   25% in 0.3479 secs
//   50% in 0.3537 secs
//   75% in 0.7003 secs
//   90% in 0.7213 secs
//   95% in 0.7292 secs
//   99% in 1.8372 secs

// SummaryResponseTimesPercentiles //
type SummaryResponseTimesPercentiles struct {
	Percentile10 float64
	Percentile25 float64
	Percentile50 float64
	Percentile75 float64
	Percentile90 float64
	Percentile95 float64
	Percentile99 float64
}

// SummaryEventsAggregationsDetail //
type SummaryEventsAggregationsDetail struct {
	Slowest float64
	Average float64
	Fastest float64
}

// SummaryStats //
type SummaryStats struct {
	Aggregations               SummaryStatsAggregations
	HTTPStatusCodeDistribution map[uint64]uint64
	ResponseTimesDistribution  []SummaryStatsResponseTimesDistributionSlot
	ResponseTimesPercentiles   SummaryResponseTimesPercentiles
	EventsAggregations         map[ReqEvt]*SummaryEventsAggregationsDetail
}

// Summary //
type Summary struct {
	Timestamp time.Time
	Spec      TestSpec
	Reports   []Report
	Stats     SummaryStats
}

// ReqEvt is Request events type //
type ReqEvt string

// Request Events Constants //
const (
	ReqEvtRequestStarted        ReqEvt = "request_started"
	ReqEvtGetConnection         ReqEvt = "get_connection"
	ReqEvtGotConnection         ReqEvt = "got_connection"
	ReqEvtGotFirstResponseByte  ReqEvt = "got_first_response_byte"
	ReqEvtGot100Continue        ReqEvt = "got_100_continue"
	ReqEvtDNSStart              ReqEvt = "dns_start"
	ReqEvtDNSDone               ReqEvt = "dns_done"
	ReqEvtConnectStart          ReqEvt = "connect_start"
	ReqEvtConnectDone           ReqEvt = "connect_done"
	ReqEvtTLSHandshakeStart     ReqEvt = "tls_handshake_start"
	ReqEvtTLSHandshakeDone      ReqEvt = "tls_handshake_done"
	ReqEvtWroteHeaders          ReqEvt = "wrote_headers"
	ReqEvtWait100Continue       ReqEvt = "wait_100_continue"
	ReqEvtWroteRequest          ReqEvt = "wrote_request"
	ReqEvtGotLastResponseByte   ReqEvt = "got_last_response_byte"
	ReqEvtRequestErrorOccurred  ReqEvt = "request_error_occurred"
	ReqEvtResponseErrorOccurred ReqEvt = "response_error_occurred"
)

// ReqEvtList (best thing next to a constant array) //
func ReqEvtList() []ReqEvt {
	return []ReqEvt{
		ReqEvtRequestStarted,
		ReqEvtGetConnection,
		ReqEvtDNSStart,
		ReqEvtDNSDone,
		ReqEvtConnectStart,
		ReqEvtConnectDone,
		ReqEvtGotConnection,
		ReqEvtTLSHandshakeStart,
		ReqEvtTLSHandshakeDone,
		ReqEvtWroteHeaders,
		ReqEvtWroteRequest,
		ReqEvtGotFirstResponseByte,
		ReqEvtGot100Continue,
		ReqEvtWait100Continue,
		ReqEvtGotLastResponseByte,
		ReqEvtRequestErrorOccurred,
		ReqEvtResponseErrorOccurred,
	}
}

// ReqEvtLabels (best thing next to a constant array) //
func ReqEvtLabels(reqEvt ReqEvt) string {
	return map[ReqEvt]string{
		ReqEvtRequestStarted:        "Request Started",
		ReqEvtGetConnection:         "Get Connection",
		ReqEvtDNSStart:              "DNS Start",
		ReqEvtDNSDone:               "DNS Done",
		ReqEvtConnectStart:          "Connect Start",
		ReqEvtConnectDone:           "Connect Done",
		ReqEvtGotConnection:         "Got Connection",
		ReqEvtTLSHandshakeStart:     "TLS Handshake Start",
		ReqEvtTLSHandshakeDone:      "TLS Handshake Done",
		ReqEvtWroteHeaders:          "Wrote Headers",
		ReqEvtWroteRequest:          "Wrote Request",
		ReqEvtGotFirstResponseByte:  "Got First Response Byte",
		ReqEvtGot100Continue:        "Got 100 Continue",
		ReqEvtWait100Continue:       "Wait 100 Continue",
		ReqEvtGotLastResponseByte:   "Got Last Response Byte",
		ReqEvtRequestErrorOccurred:  "Request Error Occurred",
		ReqEvtResponseErrorOccurred: "Response Error Occurred",
	}[reqEvt]
}
