package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestServer(t *testing.T) {

	tests := []struct {
		name                 string
		path                 string
		method               string
		expectedResponse     string
		expectedResponseCode int
	}{
		{name: "GET /index", path: "/", method: "GET", expectedResponse: "404 page not found\n", expectedResponseCode: 404},
		{name: "GET /api", path: "/api/package/mypackage-123", method: "GET", expectedResponse: "Package details for: mypackage-123", expectedResponseCode: 200},
	}

	r := mux.NewRouter()
	setupHandlers(r)

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := &http.Client{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, ts.URL+tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.Do(req)
			respBody, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			if string(respBody) != tc.expectedResponse {
				t.Errorf("Expected: %s, Got: %s", tc.expectedResponse, string(respBody))
			}

			if resp.StatusCode != tc.expectedResponseCode {
				t.Errorf("Expected response status code: %d, Got: %d", tc.expectedResponseCode, resp.StatusCode)
			}

		})
	}
}
