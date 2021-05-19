package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
		{name: "GET /api", path: "/api", method: "GET", expectedResponse: "Hello, world!", expectedResponseCode: 200},
		{name: "GET /api/", path: "/api/", method: "GET", expectedResponse: "Hello, world!", expectedResponseCode: 200},
		{name: "POST /api", path: "/api", method: "POST", expectedResponse: "I got your data", expectedResponseCode: 200},
		{name: "POST /api/", path: "/api/", method: "POST", expectedResponse: "I got your data", expectedResponseCode: 200},
		{name: "PUT /api/", path: "/api/", method: "PUT", expectedResponse: "Method not allowed\n", expectedResponseCode: 405},
		{name: "GET /healthcheck", path: "/healthcheck", method: "GET", expectedResponse: "ok", expectedResponseCode: 200},
		{name: "GET /healthcheck/", path: "/healthcheck/", method: "GET", expectedResponse: "ok", expectedResponseCode: 200},
		{name: "POST /healthcheck", path: "/healthcheck", method: "POST", expectedResponse: "Method not allowed\n", expectedResponseCode: 405},
	}

	mux := http.NewServeMux()
	setupHandlers(mux)

	ts := httptest.NewServer(mux)
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
