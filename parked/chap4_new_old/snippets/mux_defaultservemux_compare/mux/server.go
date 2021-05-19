package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func packageGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Package get handler")
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Default Handler")
}

func setupHandlers(r *mux.Router) {
	r.HandleFunc("/", defaultHandler)
	r.HandleFunc("/api/package", packageGetHandler).Methods("GET")
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
