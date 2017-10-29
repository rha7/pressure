package attacker

import (
	"sync"

	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
)

func reportGatherer(
	chanReportsSink chan apptypes.Report,
	chanReportsGathererDone chan bool,
	chanReportsGathererDoneContinue chan bool,
	reports *[]apptypes.Report,
	logger *logrus.Logger,
) {
	for {
		select {
		case report := <-chanReportsSink:
			if logger.Level == logrus.DebugLevel {
				logger.
					WithField("thread_id", report.ThreadID).
					WithField("request_id", report.RequestID).
					WithField("error", report.Error).
					WithField("outcome", apptypes.OutcomeText(report.Outcome)).
					WithField("code", report.Code).
					WithField("uncompressed", report.Uncompressed).
					WithField("timings", report.Timings).
					Info("read result")
			} else {
				logger.
					WithField("thread_id", report.ThreadID).
					WithField("request_id", report.RequestID).
					WithField("error", report.Error).
					WithField("outcome", apptypes.OutcomeText(report.Outcome)).
					WithField("code", report.Code).
					WithField("uncompressed", report.Uncompressed).
					WithField("total_time", report.Timings[len(report.Timings)-1].TimestampMilliseconds).
					Info("read result")
			}
			(*reports) = append(*reports, report)
		case <-chanReportsGathererDone:
			chanReportsGathererDoneContinue <- true
			return
		}
	}
}

// Perform //
func Perform(logger *logrus.Logger, spec apptypes.TestSpec) ([]apptypes.Report, error) {
	var results []apptypes.Report
	var chanRequestIDProvider = make(chan uint64, spec.TotalRequests)
	var chanReportsSink = make(chan apptypes.Report, spec.TotalRequests)
	var chanReportsGathererDone = make(chan bool, 1)
	var chanReportsGathererDoneContinue = make(chan bool, 1)
	waitGroup := &sync.WaitGroup{}
	go reportGatherer(
		chanReportsSink,
		chanReportsGathererDone,
		chanReportsGathererDoneContinue,
		&results,
		logger,
	)
	logger.Info("feeding request id provider")
	for requestID := uint64(0); requestID < spec.TotalRequests; requestID++ {
		chanRequestIDProvider <- requestID
	}
	close(chanRequestIDProvider)
	logger.Info("setting up request processors wait count")
	waitGroup.Add(int(spec.ConcurrentThreads))
	logger.Info("launching request processors")
	for threadID := uint64(0); threadID < spec.ConcurrentThreads; threadID++ {
		logger.WithField("thread_id", threadID).Info("launching request processor")
		go processor(threadID, chanRequestIDProvider, chanReportsSink, waitGroup, spec, logger)
	}
	waitGroup.Wait()
	logger.Info("reading results")
	chanReportsGathererDone <- true
	<-chanReportsGathererDoneContinue
	close(chanReportsSink)
	logger.Info("results read completed")
	return results, nil
}
