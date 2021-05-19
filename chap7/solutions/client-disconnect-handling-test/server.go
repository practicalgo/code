package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func doSomeWork() {
	time.Sleep(5 * time.Second)
}

func handleUserAPI(w http.ResponseWriter, r *http.Request) {
	done := make(chan bool)

	log.Println("I started processing the request")

	go func() {
		doSomeWork()
		done <- true
	}()

	select {
	case <-done:
		log.Println("doSomeWork done: Continuing request processing")
	case <-r.Context().Done():
		log.Printf("Aborting request processing: %v\n", r.Context().Err())
		return
	}

	log.Println("I finished processing the request")
}

func setupHandlers(mux *http.ServeMux) {
	timeoutDuration := 10 * time.Second

	userHandler := http.HandlerFunc(handleUserAPI)
	hTimeout := http.TimeoutHandler(
		userHandler,
		timeoutDuration,
		"I ran out of time",
	)

	mux.Handle("/api/users/", hTimeout)

}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8443"
	}

	mux := http.NewServeMux()
	setupHandlers(mux)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
