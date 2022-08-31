package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func longRunningProcess(w *io.PipeWriter) {
	for i := 0; i <= 10; i++ {
		fmt.Fprintf(w, "hello")
		time.Sleep(1 * time.Second)
	}
	w.Close()
}

func main() {

	client := http.Client{}
	logReader, logWriter := io.Pipe()
	go longRunningProcess(logWriter)
	r, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		"http://localhost:8080/api/users/",
		logReader,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting client request")

	resp, err := client.Do(r)
	if err != nil {
		log.Fatalf("Error when sending the request: %v\n", err)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
}
