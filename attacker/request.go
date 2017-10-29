package attacker

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
)

func addEvent(timings *[]apptypes.TimingEvent, name string) {
	*timings = append(
		*timings,
		apptypes.TimingEvent{
			Name:              name,
			UnixNanoTimestamp: uint64(time.Now().UnixNano()),
		},
	)
}

func rebaseEvents(timings *[]apptypes.TimingEvent) time.Time {
	t0 := (*timings)[0].UnixNanoTimestamp
	for i := range *timings {
		(*timings)[i].TimestampMilliseconds = float64((*timings)[i].UnixNanoTimestamp-t0) / float64(1000000.0)
	}
	return time.Unix(0, int64(t0))
}

func request(threadID uint64, requestID uint64, spec apptypes.TestSpec, logger *logrus.Logger) apptypes.Report {
	timings := []apptypes.TimingEvent{}
	bodyReader := bytes.NewBufferString(spec.Data)
	req, err := http.NewRequest(spec.Method, spec.URL, bodyReader)
	if err != nil {
		return apptypes.Report{
			ThreadID:     threadID,
			RequestID:    requestID,
			Code:         0,
			Response:     "",
			Error:        err.Error(),
			Outcome:      apptypes.OutcomeError,
			Uncompressed: true,
			Timings:      timings,
		}
	}
	for headerKey, headerValue := range spec.RequestHeaders {
		req.Header.Add(headerKey, headerValue)
	}
	cli := http.Client{
		Timeout: time.Duration(spec.RequestTimeout) * time.Second,
	}
	addEvent(&timings, "request_started")
	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			addEvent(&timings, "get_connection")
		},
		GotConn: func(httptrace.GotConnInfo) {
			addEvent(&timings, "got_connection")
		},
		GotFirstResponseByte: func() {
			addEvent(&timings, "got_first_response_byte")
		},
		Got100Continue: func() {
			addEvent(&timings, "got_100_continue")
		},
		DNSStart: func(httptrace.DNSStartInfo) {
			addEvent(&timings, "dns_start")
		},
		DNSDone: func(httptrace.DNSDoneInfo) {
			addEvent(&timings, "dns_done")
		},
		ConnectStart: func(network, addr string) {
			addEvent(&timings, "connect_start")
		},
		ConnectDone: func(network, addr string, err error) {
			addEvent(&timings, "connect_done")
		},
		TLSHandshakeStart: func() {
			addEvent(&timings, "tls_handshake_start")
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			addEvent(&timings, "tls_handshake_done")
		},
		WroteHeaders: func() {
			addEvent(&timings, "wrote_headers")
		},
		Wait100Continue: func() {
			addEvent(&timings, "wait_100_continue")
		},
		WroteRequest: func(httptrace.WroteRequestInfo) {
			addEvent(&timings, "wrote_request")
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := cli.Do(req)
	addEvent(&timings, "request_completed")
	timeStamp := rebaseEvents(&timings)
	if err != nil {
		return apptypes.Report{
			ThreadID:     threadID,
			RequestID:    requestID,
			Code:         0,
			Response:     "",
			Error:        err.Error(),
			Outcome:      apptypes.OutcomeError,
			Uncompressed: resp.Uncompressed,
			Timings:      timings,
			Timestamp:    timeStamp,
		}
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	addEvent(&timings, "got_last_response_byte")
	rebaseEvents(&timings)
	if err != nil {
		return apptypes.Report{
			ThreadID:     threadID,
			RequestID:    requestID,
			Code:         uint64(resp.StatusCode),
			Response:     string(bodyBytes),
			Error:        err.Error(),
			Outcome:      apptypes.OutcomeError,
			Uncompressed: resp.Uncompressed,
			Timings:      timings,
			Timestamp:    timeStamp,
		}
	}
	return apptypes.Report{
		ThreadID:     threadID,
		RequestID:    requestID,
		Code:         uint64(resp.StatusCode),
		Response:     string(bodyBytes),
		Error:        "",
		Outcome:      apptypes.OutcomeSuccess,
		Uncompressed: resp.Uncompressed,
		Timings:      timings,
		Timestamp:    timeStamp,
	}
}
