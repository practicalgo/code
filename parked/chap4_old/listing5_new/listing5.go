package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	wrappedMux := setupHandlers(mux, logger)

	log.Printf("Server attempting to listen on: %s\n", listenAddr)
	err := http.ListenAndServe(listenAddr, wrappedMux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}

func setupHandlers(mux *http.ServeMux, logger *log.Logger) http.Handler {
	mux.HandleFunc("/healthcheck", healthCheckHandler)
	mux.HandleFunc("/deepcheck", deepCheckHandler)
	mux.HandleFunc("/", defaultHandler)

	return logRequestsMiddleware(
		logger,
		panicRecoveryMiddleware(logger, mux))

}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	_, ok := req.URL.Query()["panic"]
	if ok {
		panic("Sorry, I couldn't process your request this time")
	}
	io.WriteString(w, "Hello, world!")
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "ok")
}

func deepCheckHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "deepcheck_ok")
}

func logRequestsMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("method=%s path=%s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func panicRecoveryMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				logger.Printf("stacktrace=%s\n", string(debug.Stack()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
