// pkg registration with form data
package pkgregister

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type pkgData struct {
	Name     string
	Version  string
	Filename string
	Bytes    io.Reader
}

type pkgRegisterResult struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func registerPackageData(url string, data pkgData) (pkgRegisterResult, error) {

	p := pkgRegisterResult{}
	payload, contentType, err := createMultiPartMessage(data)
	if err != nil {
		return p, err
	}
	reader := bytes.NewReader(payload)
	r, err := http.Post(url, contentType, reader)
	if err != nil {
		return p, err
	}
	defer r.Body.Close()
	respData, err := io.ReadAll(r.Body)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(respData, &p)
	return p, err
}
