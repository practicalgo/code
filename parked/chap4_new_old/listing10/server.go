package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func packageGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Package details for: %s", vars["id"])
}

func packageGetVersionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Package details for: %s Version: %s", vars["id"], vars["version"])
}

func setupHandlers(r *mux.Router) {
	r.HandleFunc("/api/package/{id:[a-z]+}", packageGetHandler).Methods("GET")
	r.HandleFunc("/api/package/{id:[a-z]+}/{version:[0-9.]+}", packageGetVersionHandler).Methods("GET")
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
