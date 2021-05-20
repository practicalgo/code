package telemetry

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
)

type DurationMetric struct {
	Cmd        string
	DurationMs float64
	Success    bool
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

func (m MetricReporter) ReportDuration(metric DurationMetric) {
	metricName := "cmd.duration"
	m.statsd.Histogram(
		metricName,
		metric.DurationMs,
		[]string{
			fmt.Sprintf("cmd=%s", metric.Cmd),
			fmt.Sprintf("success=%v", metric.Success),
		},
		1, //sample rate (0-none, 1 - all)
	)
}

func (m MetricReporter) Close() {
	m.statsd.Close()
}
