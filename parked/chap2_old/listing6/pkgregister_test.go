package pkgregister

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func packageRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Package registration response
		d := pkgRegisterResult{}
		err := r.ParseMultipartForm(5000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			mForm := r.MultipartForm
			// Get file data
			f := mForm.File["filedata"][0]
			// Construct an artificial package ID to return
			d.Id = fmt.Sprintf("%s-%s", mForm.Value["name"][0], mForm.Value["version"][0])
			d.Filename = f.Filename
			d.Size = f.Size
			// Marshal outgoing package registration response
			jsonData, err := json.Marshal(d)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, string(jsonData))
			}
		}
	} else {
		http.Error(w, "Invalid HTTP method specified", http.StatusMethodNotAllowed)
		return
	}
}

func startTestPackageServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/packages", packageRegHandler)
	ts := httptest.NewServer(mux)
	return ts
}

func TestRegisterPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()
	client := createHTTPClientWithTimeout(20 * time.Millisecond)

	p := pkgData{
		Name:     "mypackage",
		Version:  "0.1",
		Filename: "mypackage-0.1.tar.gz",
		Bytes:    strings.NewReader("data"),
	}
	pResult, err := registerPackageData(client, ts.URL+"/api/packages", p)
	if err != nil {
		t.Fatal(err)
	}

	if pResult.Id != "mypackage-0.1" {
		t.Errorf("Expected package ID to be mypackage-0.1, Got: %s", pResult.Id)
	}
	if pResult.Filename != "mypackage-0.1.tar.gz" {
		t.Errorf("Expected package filename to be mypackage-0.1.tar.gz, Got: %s", pResult.Filename)
	}
	if pResult.Size != 4 {
		t.Errorf("Expected package size to be 4, Got: %d", pResult.Size)
	}
}
