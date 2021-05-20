package middleware

import (
	"net/http"
	"time"

	"github.com/practicalgo/code/appendix-b/http-server/config"
)

func LoggingMiddleware(c *config.AppConfig, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Logger.Printf("Got request - headers:%#v\n", r.Header)
		startTime := time.Now()
		h.ServeHTTP(w, r)
		c.Logger.Info().Str(
			"protocol",
			r.Proto,
		).Str(
			"path",
			r.URL.Path,
		).Str(
			"method",
			r.Method,
		).Float64(
			"duration",
			time.Since(startTime).Seconds(),
		).Send()
	})
}
