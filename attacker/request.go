package attacker

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"regexp"
	"time"

	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
)

func addEvent(logger *logrus.Logger, threadID uint64, requestID uint64, timings *[]apptypes.TimingEvent, name apptypes.ReqEvt) {
	logger.
		WithField("thread_id", threadID).
		WithField("request_id", requestID).
		WithField("event_name", name).
		Debug("event received")
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

func getRequestDuration(timeStamp time.Time) float64 {
	return float64(uint64(time.Now().UnixNano()-timeStamp.UnixNano())) / float64(1000000.0)
}

func createHTTPClientTrace(logger *logrus.Logger, threadID uint64, requestID uint64, timings *[]apptypes.TimingEvent) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtGetConnection)
		},
		GotConn: func(httptrace.GotConnInfo) {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtGotConnection)
		},
		GotFirstResponseByte: func() {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtGotFirstResponseByte)
		},
		Got100Continue: func() {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtGot100Continue)
		},
		DNSStart: func(httptrace.DNSStartInfo) {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtDNSStart)
		},
		DNSDone: func(httptrace.DNSDoneInfo) {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtDNSDone)
		},
		ConnectStart: func(network, addr string) {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtConnectStart)
		},
		ConnectDone: func(network, addr string, err error) {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtConnectDone)
		},
		TLSHandshakeStart: func() {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtTLSHandshakeStart)
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtTLSHandshakeDone)
		},
		WroteHeaders: func() {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtWroteHeaders)
		},
		Wait100Continue: func() {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtWait100Continue)
		},
		WroteRequest: func(httptrace.WroteRequestInfo) {
			addEvent(logger, threadID, requestID, timings, apptypes.ReqEvtWroteRequest)
		},
	}
}

func request(threadID uint64, requestID uint64, client *http.Client, spec apptypes.TestSpec, logger *logrus.Logger) apptypes.Report {
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
	rxHdr := regexp.MustCompile("(?i)host")
	for headerKey, headerValue := range spec.RequestHeaders {
		if rxHdr.MatchString(headerKey) {
			req.Host = headerValue
			logger.WithField("host", headerValue).Info("setting request host")
			continue
		}
		req.Header.Add(headerKey, headerValue)
	}
	tracers := createHTTPClientTrace(logger, threadID, requestID, &timings)
	ctx := req.Context()
	ctxWithTrace := httptrace.WithClientTrace(ctx, tracers)
	req = req.WithContext(ctxWithTrace)
	resp, err := client.Do(req)
	if err != nil {
		addEvent(logger, threadID, requestID, &timings, apptypes.ReqEvtRequestErrorOccurred)
		timeStamp := rebaseEvents(&timings)
		return apptypes.Report{
			ThreadID:             threadID,
			RequestID:            requestID,
			Code:                 0,
			Response:             "ERROR: error executing request",
			Error:                err.Error(),
			Outcome:              apptypes.OutcomeError,
			Compressed:           false,
			Timings:              timings,
			Timestamp:            timeStamp,
			DurationMilliseconds: getRequestDuration(timeStamp),
		}
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		addEvent(logger, threadID, requestID, &timings, apptypes.ReqEvtResponseErrorOccurred)
		timeStamp := rebaseEvents(&timings)
		return apptypes.Report{
			ThreadID:             threadID,
			RequestID:            requestID,
			Code:                 uint64(resp.StatusCode),
			Response:             "ERROR: error reading response body",
			Error:                err.Error(),
			Outcome:              apptypes.OutcomeError,
			Compressed:           !resp.Uncompressed,
			Timings:              timings,
			Timestamp:            timeStamp,
			DurationMilliseconds: getRequestDuration(timeStamp),
		}
	}
	addEvent(logger, threadID, requestID, &timings, apptypes.ReqEvtGotLastResponseByte)
	timeStamp := rebaseEvents(&timings)
	return apptypes.Report{
		ThreadID:             threadID,
		RequestID:            requestID,
		Code:                 uint64(resp.StatusCode),
		Response:             string(bodyBytes),
		Error:                "",
		Outcome:              apptypes.OutcomeSuccess,
		Compressed:           !resp.Uncompressed,
		Timings:              timings,
		Timestamp:            timeStamp,
		DurationMilliseconds: getRequestDuration(timeStamp),
	}
}
