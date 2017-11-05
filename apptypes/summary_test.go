package apptypes_test

import (
	"testing"
	"time"

	"github.com/rha7/pressure/apptypes"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTimingEvent(t *testing.T) {
	Convey("Given a timing event", t, func() {
		location, _ := time.LoadLocation("UTC")
		timestampMilliseconds := 777.777
		unixNanoTimestamp := uint64(time.Date(1974, 7, 16, 7, 7, 7, 0, location).UnixNano())
		te := apptypes.TimingEvent{
			Name:                  apptypes.ReqEvtRequestStarted,
			UnixNanoTimestamp:     unixNanoTimestamp,
			TimestampMilliseconds: timestampMilliseconds,
		}
		Convey("When requesting string representation", func() {
			tes := te.String()
			Convey("Then I should get a string representation of it", func() {
				s := "{Name: request_started, UnixNanoTimeStamp: 143190427000000000, " +
					"TimestampMilliseconds: 777.777000}"
				So(tes, ShouldContainSubstring, s)
			})
		})
	})
}

func TestOutcomeText(t *testing.T) {
	Convey("Given an outcome", t, func() {
		outcomeSuccess := apptypes.OutcomeSuccess
		outcomeError := apptypes.OutcomeError
		outcomeUnknown := apptypes.Outcome(100)
		Convey("When converted to text", func() {
			textSuccess := apptypes.OutcomeText(outcomeSuccess)
			textError := apptypes.OutcomeText(outcomeError)
			textUnknown := apptypes.OutcomeText(outcomeUnknown)
			Convey("We should get a corresponding text representation", func() {
				So(textSuccess, ShouldEqual, "Outcome Success")
				So(textError, ShouldEqual, "Outcome Error")
				So(textUnknown, ShouldEqual, "Outcome Unknown")
			})
		})
	})
}

func TestReqEvtListAndLabels(t *testing.T) {
	Convey("Given a request/response event list", t, func() {
		evtList := apptypes.ReqEvtList()
		Convey("When checking its elements", func() {
			Convey("Then I should get the corresponding label", func() {
				So(evtList, ShouldNotBeEmpty)
				So(apptypes.ReqEvtGotFirstResponseByte, ShouldBeIn, evtList)
				So(apptypes.ReqEvtGotLastResponseByte, ShouldBeIn, evtList)
			})
		})
	})
	Convey("Given a request/response event", t, func() {
		evt := apptypes.ReqEvtRequestStarted
		Convey("When converting it to a label", func() {
			label := apptypes.ReqEvtLabels(evt)
			Convey("Then I should get the corresponding label", func() {
				So(label, ShouldEqual, "Request Started")
			})
		})
	})
}
