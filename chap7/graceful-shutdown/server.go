package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleUserAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("I started processing the request")
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v\n", err)
		http.Error(
			w, "Error reading body",
			http.StatusInternalServerError,
		)
		return
	}
	log.Println(string(data))
	fmt.Fprintf(w, "Hello world!")
	log.Println("I finished processing the request")
}

func shutDown(ctx context.Context, s *http.Server, waitForShutdownCompletion chan struct{}) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigch
	log.Printf("Got signal: %v . Server shutting down.", sig)
	if err := s.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	waitForShutdownCompletion <- struct{}{}
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	waitForShutdownCompletion := make(chan struct{})
	ctx, cancel := context.WithTimeout(
		context.Background(), 30*time.Second,
	)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/users/", handleUserAPI)

	s := http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	go shutDown(ctx, &s, waitForShutdownCompletion)

	err := s.ListenAndServe()
	log.Print(
		"Waiting for shutdown to complete..",
	)
	<-waitForShutdownCompletion
	log.Fatal(err)
}
