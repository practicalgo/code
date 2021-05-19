package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
)

func getData(ctx context.Context, w http.ResponseWriter) {
	client := http.Client{}

	r, err := http.NewRequestWithContext(ctx, "GET", "https://godoc.org", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := client.Do(r)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	io.Copy(w, resp.Body)

}

func apiHandlerFunction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	getData(ctx, w)
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
