package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func clientDisconnected(ctx context.Context, done chan bool) bool {
	select {
	case <-done:
		return false
	case <-ctx.Done():
		log.Printf("api: client disconnected.")
		return true
	}
}

func apiHandlerFunction(w http.ResponseWriter, r *http.Request) {
	done := make(chan bool)
	go func() {
		log.Println("First expensive operation")
		time.Sleep(5 * time.Second)
		done <- true
	}()

	if clientDisconnected(r.Context(), done) {
		return
	}

	go func() {
		log.Println("Second expensive operation")
		time.Sleep(5 * time.Second)
		done <- true
	}()

	if clientDisconnected(r.Context(), done) {
		return
	}

	fmt.Fprintf(w, "All operations done")
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api", apiHandlerFunction)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}

}
