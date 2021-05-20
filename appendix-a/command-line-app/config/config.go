package config

import (
	"github.com/practicalgo/code/appendix-a/pkgcli/telemetry"
	"github.com/rs/zerolog"
)

type PkgCliConfig struct {
	Logger  zerolog.Logger
	Metrics telemetry.MetricReporter
	Tracer  telemetry.TraceReporter
}
