package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func apiGetHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func apiPostHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "I got your data")
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")
}

func setupHandlers(r *mux.Router) {
	r.HandleFunc("/healthcheck", healthCheckHandler).Methods("GET")
	r.HandleFunc("/healthcheck/", healthCheckHandler).Methods("GET")

	r.HandleFunc("/api", apiGetHandler).Methods("GET")
	r.HandleFunc("/api/", apiGetHandler).Methods("GET")

	r.HandleFunc("/api", apiPostHandler).Methods("POST")
	r.HandleFunc("/api/", apiPostHandler).Methods("POST")
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	r := mux.NewRouter()
	setupHandlers(r)

	log.Printf("Server attempting to listen on: %s\n", listenAddr)
	err := http.ListenAndServe(listenAddr, r)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}
