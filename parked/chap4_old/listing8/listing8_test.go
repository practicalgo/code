package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {

	tests := []struct {
		name               string
		path               string
		expectedResponse   string
		expectedRequestLog string
		expectedHandlerLog string
	}{
		{
			name:               "index",
			path:               "/",
			expectedResponse:   "Hello, world!",
			expectedRequestLog: "method=GET path=/ status=200",
		},
		{
			name:               "index_panic",
			path:               "/?panic",
			expectedResponse:   "Sorry, I couldn't process your request this time\n",
			expectedRequestLog: "stacktrace=",
		},
		{
			name:               "index_panic",
			path:               "/?panic",
			expectedResponse:   "Sorry, I couldn't process your request this time\n",
			expectedRequestLog: "method=GET path=/?panic status=500",
		},
		{
			name:               "healthcheck",
			path:               "/healthcheck",
			expectedResponse:   "ok",
			expectedRequestLog: "method=GET path=/healthcheck status=200",
		},
		{
			name:               "healthcheck",
			path:               "/healthcheck/",
			expectedResponse:   "ok",
			expectedRequestLog: "method=GET path=/healthcheck status=200",
		},
		{
			name:               "deepcheck",
			path:               "/deepcheck",
			expectedResponse:   "deepcheck_ok",
			expectedRequestLog: "method=GET path=/deepcheck status=200",
			expectedHandlerLog: "Handling deepcheck",
		},
	}

	mux := http.NewServeMux()

	var str bytes.Buffer
	testLogger := log.New(&str, "", log.Ldate|log.Ltime|log.Lshortfile)
	wrappedMux := setupHandlers(mux, testLogger)

	var expectedLogLines []string
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()
			wrappedMux.ServeHTTP(w, req)
			res := w.Result()

			respBody, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if string(respBody) != tc.expectedResponse {
				t.Errorf("Expected Response: %s, Got: %s", tc.expectedResponse, string(respBody))
			}
		})

		expectedLogLines = append(expectedLogLines, tc.expectedRequestLog)
		if len(tc.expectedHandlerLog) != 0 {
			expectedLogLines = append(expectedLogLines, tc.expectedHandlerLog)
		}

	}
	for _, l := range expectedLogLines {
		if !strings.Contains(str.String(), l) {
			t.Fatalf("Expected Log Line: %s, Got: %s", l, str.String())
		}
	}

}
