package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestServer(t *testing.T) {

	tests := []struct {
		name             string
		path             string
		expectedResponse string
	}{
		{name: "api", path: "/api", expectedResponse: "Hello, world!"},
		{name: "healthcheck", path: "/healthcheck", expectedResponse: "ok"},
	}

	r := mux.NewRouter()
	setupHandlers(r)

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tc.path)
			respBody, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			if string(respBody) != tc.expectedResponse {
				t.Errorf("Expected: %s, Got: %s", tc.expectedResponse, string(respBody))
			}
		})
	}
}
