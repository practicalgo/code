package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func longProcessing(ctx context.Context, w http.ResponseWriter) {

	done := make(chan bool)
	go func() {
		for i := 1; i <= 100; i++ {
			time.Sleep(5 * time.Second)
			log.Printf("I am still up")
			select {
			case <-ctx.Done():
				log.Printf("Client disconnected. Cancelling expensive operation.")
				return
			default:
				continue
			}

		}
		done <- true

	}()
	<-done
}

func expensiveHandlerFunction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Running an expensive operation.")
	longProcessing(ctx, w)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/expensive", expensiveHandlerFunction)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}

}
