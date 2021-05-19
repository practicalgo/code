package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetupServer(t *testing.T) {
	b := new(bytes.Buffer)
	mux := http.NewServeMux()
	wrappedMux := setupServer(mux, b)
	ts := httptest.NewServer(wrappedMux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/panic")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf(
			"Expected response status to be: %v, Got: %v",
			http.StatusInternalServerError,
			resp.StatusCode,
		)
	}

	logs := b.String()
	expectedLogFragments := []string{
		"path=/panic method=GET duration=",
		"panic detected",
	}
	for _, log := range expectedLogFragments {
		if !strings.Contains(logs, log) {
			t.Errorf(
				"Expected logs to contain: %s, Got: %s",
				log, logs,
			)
		}
	}
}
