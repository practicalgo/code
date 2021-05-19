package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/users/", handleUserAPI)

	s := http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
