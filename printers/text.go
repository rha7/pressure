package printers

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
)

const textSummaryTemplate = `
Load Test Started At: {{.Timestamp}}

Summary:
  Requests Count: {{.Summary.Count | format_stat_uint64}}
          Errors: {{.Summary.Errors | format_stat_uint64}}
      Errors (%): {{.Summary.ErrorsPercent | format_stat_float}} %
      Total Time: {{.Summary.Total | format_stat_float}} ms
         Slowest: {{.Summary.Slowest | format_stat_float}} ms
         Fastest: {{.Summary.Fastest | format_stat_float}} ms
         Average: {{.Summary.Average | format_stat_float}} ms
    Requests/sec: {{.Summary.RequestsPerSec | format_stat_float}}

Return HTTP Status Codes Distribution:
{{range $code, $count := .CodeDistribution}}
    [{{$code}}] {{$count}} responses
{{end}}
`

// TODO: Move this to report, it belongs there, here we should only print it.

type textSummaryDataSummary struct {
	Count          uint64
	Errors         uint64
	ErrorsPercent  float64
	Total          float64
	Slowest        float64
	Fastest        float64
	Average        float64
	RequestsPerSec float64
}

type textSummaryData struct {
	Timestamp        time.Time
	Summary          textSummaryDataSummary
	CodeDistribution map[uint64]uint64
}

// TODO: Add Response time histogram (Ten slices between slowest and fastest, with count of each)
// TODO: Latency distribution (response time percentils, 10, 25, 50, 75, 90, 95, and 99)

// Text //
func Text(logger *logrus.Logger, summary apptypes.Summary) (string, error) {
	data := textSummaryData{
		Timestamp:        summary.Timestamp,
		CodeDistribution: make(map[uint64]uint64),
	}

	responseTimeOfFirst := float64(summary.Reports[0].DurationMilliseconds)
	data.Summary.Count = uint64(len(summary.Reports))
	data.Summary.Slowest = responseTimeOfFirst
	data.Summary.Fastest = responseTimeOfFirst
	for _, report := range summary.Reports {
		reportResponseTime := float64(report.DurationMilliseconds)
		if reportResponseTime > data.Summary.Slowest {
			data.Summary.Slowest = reportResponseTime
		}
		if reportResponseTime < data.Summary.Fastest {
			data.Summary.Fastest = reportResponseTime
		}
		if report.Outcome == apptypes.OutcomeError && reportResponseTime > 0.0 {
			data.Summary.Errors++
		}
		data.Summary.Total += reportResponseTime
		data.CodeDistribution[report.Code]++
	}
	data.Summary.Average = float64(data.Summary.Total / float64(data.Summary.Count))
	data.Summary.RequestsPerSec = float64(data.Summary.Count) / (data.Summary.Total / float64(1000.0))
	data.Summary.ErrorsPercent = float64(100*data.Summary.Errors) / float64(data.Summary.Count)

	funcMap := template.FuncMap{
		"format_stat_float": func(v float64) string {
			return fmt.Sprintf("%20s", humanize.FormatFloat("#,###.###", v))
		},
		"format_stat_uint64": func(v uint64) string {
			return fmt.Sprintf("%20s", humanize.Comma(int64(v)))
		},
	}

	// Proceed to printing
	b := bytes.NewBufferString("")
	t := template.Must(template.New("text_summary").Funcs(funcMap).Parse(textSummaryTemplate))
	err := t.Execute(b, data)
	if err != nil {
		return "", fmt.Errorf("error occurred while printing text: %s", err.Error())
	}
	return b.String(), nil
}
