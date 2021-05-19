package pkgquery

import (
	"encoding/json"
	"io/ioutil"
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
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return packages, err
	}
	err = json.Unmarshal(data, &packages)
	return packages, err
}

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}
