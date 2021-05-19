package middleware

import "net/http"

type customResponseWriter struct {
	http.ResponseWriter
	code int
}

func (mw *customResponseWriter) WriteHeader(code int) {
	mw.code = code
	mw.ResponseWriter.WriteHeader(code)
}
