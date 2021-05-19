package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func packageGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Package get handler")
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Default Handler")
}

func setupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/api/package", packageGetHandler)
	mux.HandleFunc("/", defaultHandler)
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	setupHandlers(mux)

	log.Printf("Server attempting to listen on: %s\n", listenAddr)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}
