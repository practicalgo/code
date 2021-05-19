package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func logCreator(logWriter *io.PipeWriter) {
	var i int
	for {
		i++
		fmt.Fprintf(logWriter, `{"id": %d, "user_ip": "172.121.19.21", "event": "click_on_add_cart" }`, i)
		fmt.Fprintln(logWriter)
		time.Sleep(1 * time.Second)
	}

}

func logStreamHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	logReader, logWriter := io.Pipe()
	defer logReader.Close()
	defer logWriter.Close()

	buf := make([]byte, 500)

	go func() {
		for {
			n, err := logReader.Read(buf)
			if err == io.EOF {
				break
			}
			w.Write(buf[:n])
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}()

	logCreator(logWriter)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/log", logStreamHandler)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}

}
