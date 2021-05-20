package telemetry

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
)

type DurationMetric struct {
	Method     string
	Success    bool
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

func (m MetricReporter) ReportLatency(name string, metric DurationMetric) {
	m.statsd.Histogram(
		name,
		metric.DurationMs,
		[]string{
			fmt.Sprintf("method:%s", metric.Method),
			fmt.Sprintf("success:%v", metric.Success),
		},
		1, //sample rate (0-none, 1 - all)
	)
}
