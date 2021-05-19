package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/practicalgolang/code/chap2/listing15/handlers"
)

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := chi.NewRouter()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	wrappedMux := handlers.SetupHandlers(mux, logger)

	log.Printf("Server attempting to listen on: %s\n", listenAddr)
	err := http.ListenAndServe(listenAddr, wrappedMux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}
