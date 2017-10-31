package attacker

import (
	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
)

func reportGatherer(
	chanReportsSink chan apptypes.Report,
	chanReportsGathererDone chan bool,
	chanReportsGathererDoneContinue chan bool,
	summary *apptypes.Summary,
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
					WithField("compressed", report.Compressed).
					WithField("total_time", report.DurationMilliseconds).
					WithField("timings", report.Timings).
					Info("read result (debug)")
			} else {
				logger.
					WithField("thread_id", report.ThreadID).
					WithField("request_id", report.RequestID).
					WithField("error", report.Error).
					WithField("outcome", apptypes.OutcomeText(report.Outcome)).
					WithField("code", report.Code).
					WithField("compressed", report.Compressed).
					WithField("total_time", report.DurationMilliseconds).
					Info("read result")
			}
			(*summary).Reports = append((*summary).Reports, report)
		case <-chanReportsGathererDone:
			chanReportsGathererDoneContinue <- true
			return
		}
	}
}
