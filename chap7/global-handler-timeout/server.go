package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handleUserAPI(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	fmt.Fprintf(w, "Hello world!")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(6 * time.Second)
	fmt.Fprintf(w, "ok")
}

func main() {

	timeoutDuration := 5 * time.Second

	mux := http.NewServeMux()
	mux.HandleFunc("/api/users/", handleUserAPI)
	mux.HandleFunc("/healthz", healthCheckHandler)

	muxT := http.TimeoutHandler(mux, timeoutDuration, "I ran out of time")
	log.Fatal(http.ListenAndServe(":8080", muxT))
}
