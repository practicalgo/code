package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

func TestHandlers(t *testing.T) {

	tests := []struct {
		name               string
		path               string
		headers            map[string]string
		expectedResponse   string
		expectedRequestLog string
		expectedHandlerLog string
	}{
		{
			name:               "index",
			path:               "/",
			headers:            map[string]string{"Authorization": "token foo"},
			expectedResponse:   "Hello user: 1",
			expectedRequestLog: "method=GET path=/ status=200",
		},
		{
			name:               "index",
			path:               "/",
			expectedResponse:   "Invalid user credentials\n",
			expectedRequestLog: "method=GET path=/ status=401",
		},
	}

	mux := chi.NewRouter()

	var str bytes.Buffer
	testLogger := log.New(&str, "", log.Ldate|log.Ltime|log.Lshortfile)
	wrappedMux := setupHandlers(mux, testLogger)

	var expectedLogLines []string
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			for k, v := range tc.headers {
				req.Header.Add(k, v)
			}
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
