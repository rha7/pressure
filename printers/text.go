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

Spec:
        Total Requests: {{.Spec.TotalRequests}}
    Concurrent Threads: {{.Spec.ConcurrentThreads}}
     Reuse Connections: {{.Spec.ReuseConnections}}
       Request Timeout: {{.Spec.RequestTimeout}} seconds
                Method: {{.Spec.Method}}
                   URL: {{.Spec.URL}}
       Request Headers: {{.Spec.RequestHeaders}}
                 Proxy: {{.Spec.Proxy}}
                  Data: {{.Spec.Data}}

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

Latency Distribution:
    10% in {{.Stats.ResponseTimesPercentiles.Percentile10 | format_stat_float}} ms
    25% in {{.Stats.ResponseTimesPercentiles.Percentile25 | format_stat_float}} ms
    50% in {{.Stats.ResponseTimesPercentiles.Percentile50 | format_stat_float}} ms
    75% in {{.Stats.ResponseTimesPercentiles.Percentile75 | format_stat_float}} ms
    90% in {{.Stats.ResponseTimesPercentiles.Percentile90 | format_stat_float}} ms
    95% in {{.Stats.ResponseTimesPercentiles.Percentile95 | format_stat_float}} ms
    99% in {{.Stats.ResponseTimesPercentiles.Percentile99 | format_stat_float}} ms

Request/Response Events Details:
{{$evtaggs := .Stats.EventsAggregations -}}
{{"                                                 Slowest                 Average                 Fastest"}}
{{"                                  ======================================================================"}}
{{range $event := event_list -}}
{{with index $evtaggs $event -}}
{{$event | event_labels | printf "%30s"}} : {{.Slowest | format_stat_float}} ms {{.Average | format_stat_float}} ms {{.Fastest | format_stat_float}} ms
{{else -}}
{{$event | event_labels | printf "%30s"}} :                    - ms                    - ms                    - ms
{{end -}}
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
		"event_list":   apptypes.ReqEvtList,
		"event_labels": apptypes.ReqEvtLabels,
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
