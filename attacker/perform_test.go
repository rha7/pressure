package attacker_test

import (
	"testing"

	"github.com/rha7/pressure/apptypes"
	"github.com/rha7/pressure/attacker"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPerform(t *testing.T) {
	Convey("Given a test spec and a logger", t, func() {
		spec := &apptypes.TestSpec{
			ConcurrentThreads: 3,
			Data:              "",
			Method:            "GET",
			Proxy:             "",
			RequestHeaders:    apptypes.Headers{},
			RequestTimeout:    10,
			ReuseConnections:  false,
			TotalRequests:     30,
			URL:               "http://www.google.com.mx/",
		}
		logger := &logrus.Logger{}
		Convey("When test is performed", func() {
			summary, err := attacker.Perform(logger, *spec)
			Convey("Then I should get no errors and a summary", func() {
				So(err, ShouldBeNil)
				So(summary, ShouldNotBeNil)
				So(summary.Reports, ShouldHaveLength, 30)
				So(summary.Stats.Aggregations.Count, ShouldEqual, 30)
			})
		})
	})
}
