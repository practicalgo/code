// pkg registration with form data
package pkgregister

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type pkgData struct {
	Name     string
	Version  string
	Filename string
	Bytes    io.Reader
}

type pkgRegisterResult struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func registerPackageData(client *http.Client, url string, data pkgData) (pkgRegisterResult, error) {

	p := pkgRegisterResult{}
	payload, contentType, err := createMultiPartMessage(data)
	if err != nil {
		return p, err
	}
	//fmt.Println(string(payload))
	reader := bytes.NewReader(payload)
	r, err := client.Post(url, contentType, reader)
	if err != nil {
		return p, err
	}
	defer r.Body.Close()
	respData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(respData, &p)
	return p, err
}

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}
