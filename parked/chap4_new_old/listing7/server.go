package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func apiGetHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func apiPostHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "I got your data")
}

func apiHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		apiGetHandler(w, req)
	case "POST":
		apiPostHandler(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		fmt.Fprintf(w, "ok")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func setupHandlers(mux *http.ServeMux) {

	mux.HandleFunc("/healthcheck", healthCheckHandler)
	mux.HandleFunc("/healthcheck/", healthCheckHandler)

	mux.HandleFunc("/api", apiHandler)
	mux.HandleFunc("/api/", apiHandler)
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	setupHandlers(mux)

	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}
