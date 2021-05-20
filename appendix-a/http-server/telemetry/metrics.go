package telemetry

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
)

type DurationMetric struct {
	Path       string
	Method     string
	DurationMs float64
}

type MetricReporter struct {
	statsd *statsd.Client
}

func InitMetrics(statsdAddr string) (MetricReporter, error) {
	var m MetricReporter
	var err error
	m.statsd, err = statsd.New(statsdAddr)
	if err != nil {
		return m, err
	}
	return m, nil
}

func (m MetricReporter) ReportLatency(metric DurationMetric) {
	metricName := "pkgserver.http.request_latency"
	m.statsd.Histogram(
		metricName,
		float64(metric.DurationMs),
		[]string{
			fmt.Sprintf("path:%s", metric.Path),
			fmt.Sprintf("method:%v", metric.Method),
		},
		1, //sample rate (0-none, 1 - all)
	)
}
