package telemetry

import (
	"io"

	"github.com/rs/zerolog"
)

func InitLogging(w io.Writer, version string, logLevel int) zerolog.Logger {

	rLogger := zerolog.New(w)
	versionedL := rLogger.With().Str("version", version)
	timestampedL := versionedL.Timestamp().Logger()
	levelledL := timestampedL.Level(zerolog.Level(logLevel))

	return levelledL
}
