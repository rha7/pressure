package attacker

import (
	"net/http"
	"sync"

	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
)

func processor(threadID uint64, chanRequestIDProvider <-chan uint64, chanResultSink chan<- apptypes.Report, client *http.Client, wg *sync.WaitGroup, spec apptypes.TestSpec, logger *logrus.Logger) {
	defer wg.Done()
	for {
		logger.
			WithField("thread_id", threadID).
			Info("waiting for request_id from provider")
		requestID, more := <-chanRequestIDProvider
		if !more {
			logger.
				WithField("thread_id", threadID).
				WithField("request_id", requestID).
				Info("exiting processor, no more request IDs available")
			break
		}
		logger.
			WithField("thread_id", threadID).
			WithField("request_id", requestID).
			Info("making request")
		result := request(threadID, requestID, client, spec, logger)
		logger.
			WithField("thread_id", threadID).
			WithField("request_id", requestID).
			Info("pushing results")
		chanResultSink <- result
		logger.
			WithField("thread_id", threadID).
			WithField("request_id", requestID).
			Info("request completed")
	}
}
