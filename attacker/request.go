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

func addEvent(timings *[]apptypes.TimingEvent, name apptypes.ReqEvt) {
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
			ThreadID:   threadID,
			RequestID:  requestID,
			Code:       0,
			Response:   "ERROR: error creating request",
			Error:      err.Error(),
			Outcome:    apptypes.OutcomeError,
			Compressed: false,
			Timings:    timings,
		}
	}
	for headerKey, headerValue := range spec.RequestHeaders {
		req.Header.Add(headerKey, headerValue)
	}
	cli := http.Client{
		Timeout: time.Duration(spec.RequestTimeout) * time.Second,
	}
	addEvent(&timings, apptypes.ReqEvtRequestStarted)
	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			addEvent(&timings, apptypes.ReqEvtGetConnection)
		},
		GotConn: func(httptrace.GotConnInfo) {
			addEvent(&timings, apptypes.ReqEvtGotConnection)
		},
		GotFirstResponseByte: func() {
			addEvent(&timings, apptypes.ReqEvtGotFirstResponseByte)
		},
		Got100Continue: func() {
			addEvent(&timings, apptypes.ReqEvtGot100Continue)
		},
		DNSStart: func(httptrace.DNSStartInfo) {
			addEvent(&timings, apptypes.ReqEvtDNSStart)
		},
		DNSDone: func(httptrace.DNSDoneInfo) {
			addEvent(&timings, apptypes.ReqEvtDNSDone)
		},
		ConnectStart: func(network, addr string) {
			addEvent(&timings, apptypes.ReqEvtConnectStart)
		},
		ConnectDone: func(network, addr string, err error) {
			addEvent(&timings, apptypes.ReqEvtConnectDone)
		},
		TLSHandshakeStart: func() {
			addEvent(&timings, apptypes.ReqEvtTLSHandshakeStart)
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			addEvent(&timings, apptypes.ReqEvtTLSHandshakeDone)
		},
		WroteHeaders: func() {
			addEvent(&timings, apptypes.ReqEvtWroteHeaders)
		},
		Wait100Continue: func() {
			addEvent(&timings, apptypes.ReqEvtWait100Continue)
		},
		WroteRequest: func(httptrace.WroteRequestInfo) {
			addEvent(&timings, apptypes.ReqEvtWroteRequest)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := cli.Do(req)
	addEvent(&timings, apptypes.ReqEvtRequestCompleted)
	timeStamp := rebaseEvents(&timings)
	if err != nil {
		return apptypes.Report{
			ThreadID:   threadID,
			RequestID:  requestID,
			Code:       0,
			Response:   "ERROR: error executing request",
			Error:      err.Error(),
			Outcome:    apptypes.OutcomeError,
			Compressed: false,
			Timings:    timings,
			Timestamp:  timeStamp,
		}
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return apptypes.Report{
			ThreadID:   threadID,
			RequestID:  requestID,
			Code:       uint64(resp.StatusCode),
			Response:   "ERROR: error reading response body",
			Error:      err.Error(),
			Outcome:    apptypes.OutcomeError,
			Compressed: !resp.Uncompressed,
			Timings:    timings,
			Timestamp:  timeStamp,
		}
	}
	addEvent(&timings, apptypes.ReqEvtGotLastResponseByte)
	rebaseEvents(&timings)
	return apptypes.Report{
		ThreadID:   threadID,
		RequestID:  requestID,
		Code:       uint64(resp.StatusCode),
		Response:   string(bodyBytes),
		Error:      "",
		Outcome:    apptypes.OutcomeSuccess,
		Compressed: !resp.Uncompressed,
		Timings:    timings,
		Timestamp:  timeStamp,
	}
}
