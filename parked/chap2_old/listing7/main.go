package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type LoggingClient struct {
	log *log.Logger
}

func (c LoggingClient) RoundTrip(r *http.Request) (*http.Response, error) {
	c.log.Printf("Sending a %s request to %s over %s\n", r.Method, r.URL, r.Proto)
	resp, err := http.DefaultTransport.RoundTrip(r)
	c.log.Printf("Got back a response over %s\n", resp.Proto)

	return resp, err
}

func main() {
	myTransport := LoggingClient{}
	l := log.New(os.Stdout, "", log.LstdFlags)
	myTransport.log = l
	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: &myTransport,
	}
	resp, err := client.Get(os.Args[1])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.Status)
	}
}
