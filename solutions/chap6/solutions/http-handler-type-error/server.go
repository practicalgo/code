package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

type appConfig struct {
	logger *log.Logger
}

type app struct {
	config  *appConfig
	handler func(rw http.ResponseWriter, req *http.Request, config *appConfig) (error, int)
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, code := a.handler(w, r, a.config)
	if err != nil {
		if len(http.StatusText(code)) == 0 {
			code = 500
			a.logger.Printf("Invalid HTTP code returned by handler while handling: %v\n", r.URL.Path)
		}
		http.Error(w, fmt.Sprintf("%s - %v", r.URL.Path, err.Error()), code)
	}
}

func apiHandler(w http.ResponseWriter, req *http.Request, config *appConfig) (error, int) {
	config.logger.Println("Handling API request")
	fmt.Fprintf(w, "Hello, world!")
	return nil, 200
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request, config *appConfig) (error, int) {
	if req.Method != "GET" {
		return errors.New("Method not allowed"), http.StatusMethodNotAllowed
	}
	config.logger.Println("Handling healthcheck request")
	fmt.Fprintf(w, "ok")
	return nil, 200
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
