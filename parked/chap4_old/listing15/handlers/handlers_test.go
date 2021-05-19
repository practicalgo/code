package handlers

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

func TestHandlers(t *testing.T) {

	tests := []struct {
		name                 string
		path                 string
		method               string
		jsonData             string
		expectedResponse     string
		expectedJsonResponse string
		expectedRequestLog   string
		expectedHandlerLog   string
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
			expectedResponse:   "Invalid login request\n",
			expectedRequestLog: "method=POST path=/ status=400",
		},
		{
			name:                 "index",
			path:                 "/",
			method:               "POST",
			jsonData:             `{"username":"joe", "password":"pass123"}`,
			expectedJsonResponse: `{"id":1,"username":"joe"}`,
			expectedRequestLog:   "method=POST path=/ status=200",
		},
		{
			name:               "index",
			path:               "/",
			method:             "POST",
			jsonData:           `{"username":"1joe", "password":"pass123"}`,
			expectedResponse:   "Invalid login data\n",
			expectedRequestLog: "method=POST path=/ status=400",
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
	wrappedMux := SetupHandlers(mux, testLogger)

	var expectedLogLines []string
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpMethod := "GET"
			if len(tc.method) != 0 {
				httpMethod = tc.method
			}
			var data io.Reader
			if len(tc.jsonData) != 0 {
				data = strings.NewReader(tc.jsonData)
			}
			req := httptest.NewRequest(httpMethod, tc.path, data)
			w := httptest.NewRecorder()
			wrappedMux.ServeHTTP(w, req)
			res := w.Result()

			respBody, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			if len(tc.expectedResponse) != 0 {
				if string(respBody) != tc.expectedResponse {
					t.Errorf("Expected Response: %s, Got: %s", tc.expectedResponse, string(respBody))
				}
			} else {
				if string(respBody) != tc.expectedJsonResponse {
					t.Errorf("Expected Response: %s, Got: %s", tc.expectedJsonResponse, string(respBody))
				}
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
