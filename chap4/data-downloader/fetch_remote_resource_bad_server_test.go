package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func startBadTestHTTPServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(60 * time.Second)
		fmt.Fprint(w, "Hello World")
	}))
	return ts
}

func TestFetchBadRemoteResource(t *testing.T) {
	ts := startBadTestHTTPServer()
	defer ts.Close()

	data, err := fetchRemoteResource(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello World"
	got := string(data)

	if expected != got {
		t.Errorf("Expected response to be: %s, Got: %s", expected, got)
	}
}
