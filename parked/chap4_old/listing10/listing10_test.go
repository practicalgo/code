package main

import (
	"bytes"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {

	tests := []struct {
		name               string
		path               string
		method             string
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
			name:               "index",
			path:               "/",
			method:             "POST",
			expectedResponse:   "Hello, world! POST handler",
			expectedRequestLog: "method=GET path=/ status=200",
		},
		{
			name:               "index",
			path:               "/",
			method:             "PUT",
			expectedResponse:   "",
			expectedRequestLog: "method=PUT path=/ status=405",
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
			name:               "admin",
			path:               "/admin/login",
			expectedResponse:   "admin_login",
			expectedRequestLog: "method=GET path=/ status=200",
		},
		{
			name:               "admin_1",
			path:               "/admin/1",
			expectedResponse:   "admin_1",
			expectedRequestLog: "method=GET path=/ status=200",
		},
		{
			name:               "admin_ana",
			path:               "/admin/ana",
			expectedResponse:   "admin_ana",
			expectedRequestLog: "method=GET path=/ status=200",
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

	mux := chi.NewRouter()
	var str bytes.Buffer
	testLogger := log.New(&str, "", log.Ldate|log.Ltime|log.Lshortfile)
	wrappedMux := setupHandlers(mux, testLogger)

	var expectedLogLines []string
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpMethod := "GET"
			if len(tc.method) != 0 {
				httpMethod = tc.method
			}
			req := httptest.NewRequest(httpMethod, tc.path, nil)
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
