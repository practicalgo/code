package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{name: "index", path: "/api", expected: "Hello, world!"},
		{name: "healthcheck", path: "/healthcheck", expected: "ok"},
	}

	setupHandlers()

	ts := httptest.NewServer(nil)
	defer ts.Close()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tc.path)
			respBody, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			if string(respBody) != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, string(respBody))
			}
		})
	}
}
