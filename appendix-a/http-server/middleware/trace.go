package middleware

import (
	"net/http"

	"github.com/practicalgo/code/appendix-a/http-server/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func TracingMiddleware(c *config.AppConfig, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Trace = otel.Tracer("")
		tc := propagation.TraceContext{}
		incomingCtx := tc.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		c.TraceCtx = incomingCtx

		ctx, span := c.Trace.Start(c.TraceCtx, r.URL.Path)
		c.Span = span
		c.SpanCtx = ctx
		defer c.Span.End()
		h.ServeHTTP(w, r)
	})
}
