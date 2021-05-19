package pkgquery

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func fetchPackageData(client *http.Client, req *http.Request) ([]pkgData, error) {
	var packages []pkgData
	r, err := client.Do(req)
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

func createHTTPGetRequest(ctx context.Context, url string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return req, err
}
