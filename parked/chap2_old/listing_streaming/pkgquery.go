package pkgquery

import (
	"encoding/json"
	"net/http"
	"time"
)

type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func fetchPackageData(client *http.Client, url string) ([]pkgData, error) {
	var packages []pkgData
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.Header.Get("Content-Type") != "application/json" {
		return packages, nil
	}
	d := json.NewDecoder(r.Body)
	err = d.Decode(&packages)
	return packages, err
}

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}
