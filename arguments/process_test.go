package arguments_test

import (
	"bytes"
	"testing"

	"github.com/rha7/pressure/arguments"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProcess(t *testing.T) {
	Convey("Given a set of arguments", t, func() {
		testURL := "http://this.issome.url/"
		args := []string{}
		bodyBuffer := bytes.NewBufferString("This is the input")
		outputBuffer := &bytes.Buffer{}
		loggerBuffer := &bytes.Buffer{}
		logger := logrus.New()
		logger.Out = loggerBuffer
		Convey("When processed with defaults and a valid URL", func() {
			args := append(args, testURL)
			testSpec, logLevel, err := arguments.Process(args, bodyBuffer, logger, outputBuffer)
			Convey("Then we should get a valid spec", func() {
				So(err, ShouldBeNil)
				So(logLevel, ShouldEqual, logrus.InfoLevel)
				So(testSpec, ShouldNotBeNil)
				So(testSpec.ConcurrentThreads, ShouldEqual, 10)
				So(testSpec.Data, ShouldEqual, "")
				So(testSpec.Method, ShouldEqual, "GET")
				So(testSpec.RequestHeaders, ShouldHaveLength, 0)
				So(testSpec.RequestTimeout, ShouldEqual, 60)
				So(testSpec.TotalRequests, ShouldEqual, 100)
				So(testSpec.URL, ShouldEqual, testURL)
				So(outputBuffer.String(), ShouldBeEmpty)
			})
		})
		Convey("When processed with an invalid flag", func() {
			args = append(args, "-invalid")
			_, _, err := arguments.Process(args, bodyBuffer, logger, outputBuffer)
			Convey("Then we should get an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "-invalid")
				So(outputBuffer.String(), ShouldContainSubstring, "Usage of")
			})
		})
		Convey("When processed with an invalid method", func() {
			args = append(args, "-X=INVALID")
			args = append(args, testURL)
			_, _, err := arguments.Process(args, bodyBuffer, logger, outputBuffer)
			Convey("Then we should get an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid method")
				So(outputBuffer.String(), ShouldBeEmpty)
			})
		})
		Convey("When processed with an invalid number of requests", func() {
			args = append(args, "-n=0")
			args = append(args, testURL)
			_, _, err := arguments.Process(args, bodyBuffer, logger, outputBuffer)
			Convey("Then we should get an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "number of requests")
				So(outputBuffer.String(), ShouldBeEmpty)
			})
		})
		Convey("When processed with an invalid log level", func() {
			args = append(args, "-l=invalid")
			_, _, err := arguments.Process(args, bodyBuffer, logger, outputBuffer)
			Convey("Then we should get an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid log level")
				So(outputBuffer.String(), ShouldBeEmpty)
			})
		})
		Convey("When processed without a valid URL", func() {
			_, _, err := arguments.Process(args, bodyBuffer, logger, outputBuffer)
			Convey("Then we should get an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "URL")
				So(outputBuffer.String(), ShouldBeEmpty)
			})
		})
	})
}
