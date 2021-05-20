package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientDisconnect(t *testing.T) {

	m := http.NewServeMux()
	setupHandlers(m)

	ts := httptest.NewServer(m)
	defer ts.Close()

	client := http.Client{
		Timeout: 100 * time.Millisecond,
	}
	resp, err := client.Get(ts.URL + "/api/users")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))

}
