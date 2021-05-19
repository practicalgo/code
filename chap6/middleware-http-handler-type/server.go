package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type appConfig struct {
	logger *log.Logger
}

type app struct {
	config  *appConfig
	handler func(w http.ResponseWriter, r *http.Request, config *appConfig)
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	a.handler(w, r, a.config)
	a.config.logger.Printf("path=%s method=%s duration=%f", r.URL.Path, r.Method, time.Now().Sub(startTime).Seconds())
}

func apiHandler(w http.ResponseWriter, r *http.Request, config *appConfig) {
	fmt.Fprintf(w, "Hello, world!")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, config *appConfig) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "ok")
}

func setupHandlers(mux *http.ServeMux, config *appConfig) {
	mux.Handle("/healthz", &app{config: config, handler: healthCheckHandler})
	mux.Handle("/api", &app{config: config, handler: apiHandler})
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	config := appConfig{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	mux := http.NewServeMux()
	setupHandlers(mux, &config)

	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
