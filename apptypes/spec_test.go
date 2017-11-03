package apptypes_test

import (
	"bytes"
	"testing"

	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHeaders(t *testing.T) {
	Convey("Given a header", t, func() {
		h := apptypes.Headers{
			"Content-Type": "application/json",
		}
		Convey("When asking for its string representation", func() {
			s := h.String()
			Convey("Then I should get it", func() {
				So(s, ShouldEqual,
					"apptypes.Headers{\"Content-Type\":\"application/json\"}")
			})
		})
		Convey("When setting a new header", func() {
			err := h.Set("Accept-Encoding: gzip")
			Convey("Then I should be able to retrieve it", func() {
				So(h["Accept-Encoding"], ShouldEqual, "gzip")
				So(err, ShouldBeNil)
			})
		})
		Convey("When setting an invalid header", func() {
			err := h.Set("Invalid header")
			Convey("Then I should get an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid header format")
			})
		})
	})
}

func TestTestSpec(t *testing.T) {
	Convey("Given a valid test spec", t, func() {
		ts := apptypes.TestSpec{
			ConcurrentThreads: 1,
			Data:              "Some Data",
			Method:            "POST",
			RequestHeaders: apptypes.Headers{
				"Content-Type":    "application/json",
				"Accept-Encoding": "gzip",
			},
			RequestTimeout: 60,
			TotalRequests:  100,
			URL:            "http://www.google.com.mx/",
		}
		Convey("When asking to log its value", func() {
			stringRecorder := &bytes.Buffer{}
			logger := logrus.New()
			logger.Level = logrus.InfoLevel
			logger.Out = stringRecorder
			ts.Print(logger)
			So(stringRecorder.String(), ShouldContainSubstring, "total_requests: 100")
			So(stringRecorder.String(), ShouldContainSubstring, "concurrent_threads: 1")
			So(stringRecorder.String(), ShouldContainSubstring, "request_timeout: 60")
			So(stringRecorder.String(), ShouldContainSubstring, "method: POST")
			So(stringRecorder.String(), ShouldContainSubstring, "url: http://www.google.com.mx/")
			So(stringRecorder.String(), ShouldContainSubstring, "request_headers:")
			So(stringRecorder.String(), ShouldContainSubstring, "Accept-Encoding: gzip")
			So(stringRecorder.String(), ShouldContainSubstring, "Content-Type: application/json")
			So(stringRecorder.String(), ShouldContainSubstring, "data: Some Data")
		})
	})
}
