package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func setupHandlersAndMiddleware(mux *http.ServeMux, l *log.Logger) http.Handler {
	mux.HandleFunc("/api", apiHandler)
	return loggingMiddleware(mux, l)
}

func loggingMiddleware(h http.Handler, l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		l.Printf("protocol=%s path=%s method=%s duration=%f", r.Proto, r.URL.Path, r.Method, time.Now().Sub(startTime).Seconds())
	})
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8443"
	}

	tlsCertFile := os.Getenv("TLS_CERT_FILE_PATH")
	tlsKeyFile := os.Getenv("TLS_KEY_FILE_PATH")

	if len(tlsCertFile) == 0 || len(tlsKeyFile) == 0 {
		log.Fatal("TLS_CERT_FILE_PATH and TLS_KEY_FILE_PATH must both be specified")
	}

	mux := http.NewServeMux()

	l := log.New(os.Stdout, "tls-server", log.Lshortfile|log.LstdFlags)
	m := setupHandlersAndMiddleware(mux, l)

	log.Fatal(http.ListenAndServeTLS(listenAddr, tlsCertFile, tlsKeyFile, m))
}
