package pkgquery

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func startTestPackageServer() *httptest.Server {
	pkgData := `[
{"name": "package1", "version": "1.1"},
{"name": "package2", "version": "1.0"}
]`
	mux := http.NewServeMux()
	mux.HandleFunc("/api/packages", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, pkgData)
	}))

	mux.HandleFunc("/packages", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "Packages available")
	}))

	ts := httptest.NewServer(mux)
	return ts
}

func TestFetchPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()
	client := createHTTPClientWithTimeout(20 * time.Millisecond)

	packages, err := fetchPackageData(client, ts.URL+"/api/packages")
	if err != nil {
		t.Fatal(err)
	}
	if len(packages) != 2 {
		t.Logf("Expected 2 packages, Got back: %d", len(packages))
		t.Fail()
	}

	packages, err = fetchPackageData(client, ts.URL+"/packages")
	if err != nil {
		t.Fatal(err)
	}
	if len(packages) != 0 {
		t.Logf("Expected 0 packages, Got back: %d", len(packages))
		t.Fail()
	}
}
