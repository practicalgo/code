package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func apiHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")
}

func setupHandlers() {
	http.HandleFunc("/healthcheck", healthCheckHandler)
	http.HandleFunc("/api", apiHandler)

}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	setupHandlers()

	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}
