package middleware

import (
	"net/http"
	"time"

	"github.com/practicalgo/code/appendix-b/http-server/config"
	"github.com/practicalgo/code/appendix-b/http-server/telemetry"
)

func MetricMiddleware(c *config.AppConfig, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(startTime).Seconds()
		c.Metrics.ReportLatency(
			telemetry.DurationMetric{
				DurationMs: duration,
				Path:       r.URL.Path,
				Method:     r.Method,
			},
		)

	})
}
