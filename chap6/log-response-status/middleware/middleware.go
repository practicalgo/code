package middleware

import (
	"net/http"
	"time"

	"github.com/practicalgo/code/chap6/log-response-status/config"
)

func loggingMiddleware(h http.Handler, config config.AppConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customRw := customResponseWriter{ResponseWriter: w, code: -1}
		startTime := time.Now()
		h.ServeHTTP(&customRw, r)
		if customRw.code == -1 {
			customRw.code = http.StatusInternalServerError
		}
		config.Logger.Printf("protocol=%s path=%s method=%s duration=%f status=%d", r.Proto, r.URL.Path, r.Method, time.Now().Sub(startTime).Seconds(), customRw.code)
	})
}

func panicMiddleware(h http.Handler, config config.AppConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rValue := recover(); rValue != nil {
				config.Logger.Println("panic detected when handling request:", rValue)
				http.Error(w, "Unexpected error", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
