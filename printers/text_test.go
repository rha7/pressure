package printers_test

import (
	"testing"

	"github.com/rha7/pressure/apptypes"
	"github.com/rha7/pressure/printers"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestText(t *testing.T) {
	Convey("Given a Text printer, a logger and a summary", t, func() {
		logger := &logrus.Logger{}
		summary := apptypes.Summary{
			Stats: apptypes.SummaryStats{
				ResponseTimesDistribution: []apptypes.SummaryStatsResponseTimesDistributionSlot{
					apptypes.SummaryStatsResponseTimesDistributionSlot{
						Count: 10,
					},
					apptypes.SummaryStatsResponseTimesDistributionSlot{
						Count: 20,
					},
				},
			},
		}
		Convey("When I use the text printer", func() {
			s, err := printers.Default(logger, summary)
			Convey("Then I should get a proper text output", func() {
				So(err, ShouldBeNil)
				So(s, ShouldContainSubstring, "Load Test Started At")
			})
		})
	})
}
