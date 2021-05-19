package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/practicalgolang/code/chap2/listing18/handlers"
)

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	// Create and setup handlers
	mux := chi.NewRouter()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	wrappedMux := handlers.SetupHandlers(mux, logger)

	// Create http server object
	srv := http.Server{
		Handler:      wrappedMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server attempting to listen on: %s\n", listenAddr)

	// Create network listener
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}
	defer l.Close()

	// Create the HTTP server
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Printf("Got SIGINT. Server shutting down..")

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.Serve(l); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
