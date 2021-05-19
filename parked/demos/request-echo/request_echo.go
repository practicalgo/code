package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Host: %s\n", req.Host)
	fmt.Fprintf(w, "Method: %s\n", req.Method)
	fmt.Fprintf(w, "Proto: %s ProtoMajor: %d ProtoManor: %d\n", req.Proto, req.ProtoMajor, req.ProtoMinor)
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header - %s:%#v\n", k, v)
	}

	fmt.Fprintf(w, "Path: %s\n", req.URL.Path)
	fmt.Fprintf(w, "Query parameters: %#v\n", req.URL.Query())
	fmt.Fprintf(w, "Username: %s\n", req.URL.User.Username())
	p, ok := req.URL.User.Password()
	if ok {
		fmt.Fprintf(w, "Password: %s \n", p)
	}

}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}
