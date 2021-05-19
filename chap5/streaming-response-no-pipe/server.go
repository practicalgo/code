package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func longRunningProcess(w http.ResponseWriter) {
	for i := 0; i <= 20; i++ {
		fmt.Fprintf(w, `{"id": %d, "user_ip": "172.121.19.21", "event": "click_on_add_cart" }`, i)
		fmt.Fprintln(w)
		f, flushSupported := w.(http.Flusher)
		if flushSupported {
			f.Flush()
		}

		time.Sleep(1 * time.Second)
	}
}

func longRunningProcessHandler(w http.ResponseWriter, r *http.Request) {
	longRunningProcess(w)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/job", longRunningProcessHandler)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}

}
