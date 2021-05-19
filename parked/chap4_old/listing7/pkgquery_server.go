package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"os"
)

type pkgData struct {
	ID      int    `json:"id" xml:"package>id"`
	Name    string `json:"name" xml:"package>name"`
	Version string `json:"version" xml:"package>version"`
}

func getPackages() []pkgData {
	return []pkgData{
		{ID: 1, Name: "package1", Version: "1.0"},
		{ID: 2, Name: "package2", Version: "1.2"},
	}
}

func pkgQueryHandler(w http.ResponseWriter, r *http.Request) {

	packages := getPackages()
	contentType := r.Header.Get("Accept")

	sendJSON := func() {
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.Encode(packages)
	}
	switch contentType {

	case "application/xml":
		w.Header().Set("Content-Type", "application/xml")
		enc := xml.NewEncoder(w)
		enc.Encode(packages)

	case "application/json":
		sendJSON()

	default:
		sendJSON()
	}
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/packages", pkgQueryHandler)

	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}
