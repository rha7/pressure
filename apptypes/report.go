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
	Name                  ReqEvt  `yaml:"name"`
	UnixNanoTimestamp     uint64  `yaml:"unix_nano_timestamp"`
	TimestampMilliseconds float64 `yaml:"timestamp_milliseconds"`
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
	ThreadID   uint64        `yaml:"thread_id"`
	RequestID  uint64        `yaml:"request_id"`
	Outcome    Outcome       `yaml:"outcome"`
	Error      string        `yaml:"error"`
	Code       uint64        `yaml:"code"`
	Response   string        `yaml:"response"`
	Compressed bool          `yaml:"compressed"`
	Timings    []TimingEvent `yaml:"timings"`
	Timestamp  time.Time     `yaml:"timestamp"`
}

// Summary //
type Summary struct {
	Timestamp time.Time `yaml:"timestamp"`
	Spec      TestSpec  `yaml:"test_spec"`
	Reports   []Report  `yaml:"reports"`
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
	}
}
