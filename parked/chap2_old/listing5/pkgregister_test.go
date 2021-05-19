package pkgregister

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func packageRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Incoming package data
		p := pkgData{}

		// Package registration response
		d := pkgRegisterResult{}
		defer r.Body.Close()
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			err := json.Unmarshal(data, &p)
			if err != nil || len(p.Name) == 0 || len(p.Version) == 0 {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			} else {
				d.Id = p.Name + "-" + p.Version
				jsonData, err := json.Marshal(d)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, string(jsonData))
				}
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
		Name:    "mypackage",
		Version: "0.1",
	}
	resp, err := registerPackageData(client, ts.URL+"/api/packages", p)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Id != "mypackage-0.1" {
		t.Errorf("Expected package id to be mypackage-0.1, Got: %s", resp.Id)
	}
}

func TestRegisterEmptyPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()
	p := pkgData{}
	client := createHTTPClientWithTimeout(20 * time.Millisecond)
	resp, err := registerPackageData(client, ts.URL+"/api/packages", p)
	if err == nil {
		t.Fatal("Expected error to be non-nil, got nil")
	}
	if len(resp.Id) != 0 {
		t.Errorf("Expected package ID to be empty, got: %s", resp.Id)
	}
}
