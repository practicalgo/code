package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{name: "index", path: "/", expected: "Hello, world!"},
		{name: "healthcheck", path: "/healthcheck", expected: "ok"},
		{name: "deepcheck", path: "/deepcheck", expected: "deepcheck_ok"},
	}

	mux := http.NewServeMux()
	setupHandlers(mux)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			res := w.Result()

			respBody, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if string(respBody) != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, string(respBody))
			}
		})
	}
}
