package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	// Check if LISTEN_ADDR environmant variable has been specified
	// if yes, use that, else default to ":8080"
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	// Setup the handlers
	setupHandlers(mux)

	log.Printf("Server attempting to listen on: %s\n", listenAddr)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}

func setupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/healthcheck", healthCheckHandler)
	mux.HandleFunc("/deepcheck", deepCheckHandler)
	mux.HandleFunc("/", defaultHandler)
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!")
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "ok")
}

func deepCheckHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "deepcheck_ok")
}
