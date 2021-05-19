package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

type requestContextKey struct{}
type requestContextValue struct {
	requestID string
}

func addRequestID(r *http.Request, requestID string) *http.Request {
	c := requestContextValue{
		requestID: requestID,
	}
	currentCtx := r.Context()
	newCtx := context.WithValue(currentCtx, requestContextKey{}, c)
	return r.WithContext(newCtx)

}

func logRequest(r *http.Request) {
	ctx := r.Context()
	v := ctx.Value(requestContextKey{})

	if m, ok := v.(requestContextValue); ok {
		log.Printf("Processing request: %s", m.requestID)
	}
}

func processRequest(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	fmt.Fprintf(w, "Request processed")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	requestID := "request-123-abc"
	r = addRequestID(r, requestID)
	processRequest(w, r)
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api", apiHandler)

	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
