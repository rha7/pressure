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

// SummaryStats //
type SummaryStats struct {
	Aggregations               SummaryStatsAggregations
	HTTPStatusCodeDistribution map[uint64]uint64
	ResponseTimesDistribution  []SummaryStatsResponseTimesDistributionSlot
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
	ReqEvtRequestStarted       ReqEvt = "request_started"
	ReqEvtGetConnection        ReqEvt = "get_connection"
	ReqEvtGotConnection        ReqEvt = "got_connection"
	ReqEvtGotFirstResponseByte ReqEvt = "got_first_response_byte"
	ReqEvtGot100Continue       ReqEvt = "got_100_continue"
	ReqEvtDNSStart             ReqEvt = "dns_start"
	ReqEvtDNSDone              ReqEvt = "dns_done"
	ReqEvtConnectStart         ReqEvt = "connect_start"
	ReqEvtConnectDone          ReqEvt = "connect_done"
	ReqEvtTLSHandshakeStart    ReqEvt = "tls_handshake_start"
	ReqEvtTLSHandshakeDone     ReqEvt = "tls_handshake_done"
	ReqEvtWroteHeaders         ReqEvt = "wrote_headers"
	ReqEvtWait100Continue      ReqEvt = "wait_100_continue"
	ReqEvtWroteRequest         ReqEvt = "wrote_request"
	ReqEvtRequestCompleted     ReqEvt = "request_completed"
	ReqEvtGotLastResponseByte  ReqEvt = "got_last_response_byte"
	ReqEvtRequestErrorOcurred  ReqEvt = "request_error_occurred"
	ReqEvtResponseErrorOcurred ReqEvt = "response_error_occurred"
)

// ReqEvtList (best thing next to a constant array) //
func ReqEvtList() []ReqEvt {
	return []ReqEvt{
		ReqEvtRequestStarted,
		ReqEvtGetConnection,
		ReqEvtGotConnection,
		ReqEvtGotFirstResponseByte,
		ReqEvtGot100Continue,
		ReqEvtDNSStart,
		ReqEvtDNSDone,
		ReqEvtConnectStart,
		ReqEvtConnectDone,
		ReqEvtTLSHandshakeStart,
		ReqEvtTLSHandshakeDone,
		ReqEvtWroteHeaders,
		ReqEvtWait100Continue,
		ReqEvtWroteRequest,
		ReqEvtRequestCompleted,
		ReqEvtGotLastResponseByte,
		ReqEvtRequestErrorOcurred,
		ReqEvtResponseErrorOcurred,
	}
}
