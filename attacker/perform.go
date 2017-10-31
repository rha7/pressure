package attacker

import (
	"sync"
	"time"

	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
)

// Perform //
func Perform(logger *logrus.Logger, spec apptypes.TestSpec) (apptypes.Summary, error) {
	var chanRequestIDProvider = make(chan uint64, spec.TotalRequests)
	var chanReportsSink = make(chan apptypes.Report, spec.TotalRequests)
	var chanReportsGathererDone = make(chan bool, 1)
	var chanReportsGathererDoneContinue = make(chan bool, 1)
	summary := apptypes.Summary{
		Timestamp: time.Now(),
		Spec:      spec,
		Reports:   []apptypes.Report{},
	}
	waitGroup := &sync.WaitGroup{}
	go reportGatherer(
		chanReportsSink,
		chanReportsGathererDone,
		chanReportsGathererDoneContinue,
		&summary,
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
	logger.Info("reading summary")
	chanReportsGathererDone <- true
	<-chanReportsGathererDoneContinue
	close(chanReportsSink)
	logger.Info("summary read completed")
	enrichSummary(&summary)
	return summary, nil
}
