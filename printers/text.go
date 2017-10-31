package printers

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/dustin/go-humanize"
	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
)

const textSummaryTemplate = `
Load Test Started At: {{.Timestamp}}

Summary:
  Requests Count: {{.Stats.Aggregations.Count | format_stat_uint64}}
          Errors: {{.Stats.Aggregations.Errors | format_stat_uint64}}
      Errors (%): {{.Stats.Aggregations.ErrorsPercent | format_stat_float}} %
      Total Time: {{.Stats.Aggregations.Total | format_stat_float}} ms
         Slowest: {{.Stats.Aggregations.Slowest | format_stat_float}} ms
         Fastest: {{.Stats.Aggregations.Fastest | format_stat_float}} ms
         Average: {{.Stats.Aggregations.Average | format_stat_float}} ms
    Requests/sec: {{.Stats.Aggregations.RequestsPerSec | format_stat_float}}

Return HTTP Status Codes Distribution:
{{range $code, $count := .Stats.HTTPStatusCodeDistribution}}
    [{{$code}}] {{$count}} responses
{{end}}

Response Times Distribution (count for greater than or equal to each slot):
{{range $code, $slot := .Stats.ResponseTimesDistribution -}}
    {{$slot.From | format_stat_float}} : {{$slot.Count | format_stat_uint64_left}} : {{$slot.Count | histogram}}
{{end}}
`

// Text //
func Text(logger *logrus.Logger, summary apptypes.Summary) (string, error) {
	maxRTD := summary.Stats.ResponseTimesDistribution[0].Count
	for _, rtd := range summary.Stats.ResponseTimesDistribution {
		if rtd.Count > maxRTD {
			maxRTD = rtd.Count
		}
	}
	funcMap := template.FuncMap{
		"format_stat_float": func(v float64) string {
			return fmt.Sprintf("%20s", humanize.FormatFloat("#,###.###", v))
		},
		"format_stat_uint64": func(v uint64) string {
			return fmt.Sprintf("%20s", humanize.Comma(int64(v)))
		},
		"format_stat_uint64_left": func(v uint64) string {
			return fmt.Sprintf("%-20s", humanize.Comma(int64(v)))
		},
		"histogram": func(count uint64) string {
			toPrint := uint64((float64(count) / float64(maxRTD)) * float64(80))
			s := ""
			for i := uint64(0); i < toPrint; i++ {
				s += "\u2588"
			}
			return s
		},
	}

	// Proceed to printing
	b := bytes.NewBufferString("")
	t := template.Must(template.New("text_summary").Funcs(funcMap).Parse(textSummaryTemplate))
	err := t.Execute(b, summary)
	if err != nil {
		return "", fmt.Errorf("error occurred while printing text: %s", err.Error())
	}
	return b.String(), nil
}
