package pkgregister

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type pkgRegisterResult struct {
	Id string `json:"id"`
}

func registerPackageData(client *http.Client, url string, data pkgData) (pkgRegisterResult, error) {
	p := pkgRegisterResult{}
	b, err := json.Marshal(data)
	if err != nil {
		return p, err
	}
	reader := bytes.NewReader(b)
	r, err := client.Post(url, "application/json", reader)
	if err != nil {
		return p, err
	}
	defer r.Body.Close()
	respData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return p, err
	}
	if r.StatusCode == 200 {
		err = json.Unmarshal(respData, &p)
	} else {
		err = errors.New(string(respData))
	}
	return p, err
}

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}
