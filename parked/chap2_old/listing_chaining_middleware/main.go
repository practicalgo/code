package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type AddHeadersMiddleware struct {
	headers map[string]string
}

type LoggingMiddleware struct {
	log *log.Logger
}

type Middleware struct {
	rtp  http.RoundTripper
	next *Middleware
}

func (m Middleware) RoundTrip(r *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	resp, err = m.rtp.RoundTrip(r)
	for {
		if m.next == nil {
			resp, err = http.DefaultTransport.RoundTrip(r)
			break
		} else {
			m.rtp.RoundTrip(r)
			m = *m.next
		}
	}
	return resp, err
}

func (c LoggingMiddleware) RoundTrip(r *http.Request) (*http.Response, error) {
	c.log.Printf("Sending a %s request to %s over %s with headers: %v\n", r.Method, r.URL, r.Proto, r.Header)
	resp, err := http.DefaultTransport.RoundTrip(r)
	c.log.Printf("Got back a response over %s\n", resp.Proto)

	return resp, err
}

func (h AddHeadersMiddleware) RoundTrip(r *http.Request) (*http.Response, error) {

	reqCopy := r.Clone(r.Context())
	for k, v := range h.headers {
		reqCopy.Header.Add(k, v)
	}
	return http.DefaultTransport.RoundTrip(reqCopy)
}

func main() {
	h := AddHeadersMiddleware{
		headers: map[string]string{
			"X-Auth-Token": "foobar",
		},
	}
	l := LoggingMiddleware{log: log.New(os.Stdout, "", log.LstdFlags)}

	mHead := Middleware{rtp: h}
	mLog := Middleware{rtp: l}
	mHead.next = &mLog

	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: mHead,
	}
	resp, err := client.Get(os.Args[1])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.Status)
	}
}
