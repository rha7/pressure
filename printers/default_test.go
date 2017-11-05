package printers_test

import (
	"testing"

	"github.com/rha7/pressure/apptypes"
	"github.com/rha7/pressure/printers"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDefault(t *testing.T) {
	Convey("Given a Default printer, a logger and a summary", t, func() {
		logger := &logrus.Logger{}
		summary := apptypes.Summary{
			Stats: apptypes.SummaryStats{
				ResponseTimesDistribution: []apptypes.SummaryStatsResponseTimesDistributionSlot{
					apptypes.SummaryStatsResponseTimesDistributionSlot{
						Count: 10,
					},
				},
			},
		}
		Convey("When I use the default printer", func() {
			s, err := printers.Default(logger, summary)
			Convey("Then I should get a proper output", func() {
				So(err, ShouldBeNil)
				So(s, ShouldNotBeEmpty)
			})
		})
	})
}
