package middleware

import (
	"net/http"

	"github.com/practicalgo/code/chap6/log-response-status/config"
)

func RegisterMiddleware(mux *http.ServeMux, conf config.AppConfig) http.Handler {
	return loggingMiddleware(panicMiddleware(mux, conf), conf)
}
