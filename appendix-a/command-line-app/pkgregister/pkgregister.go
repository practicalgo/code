// pkg registration with form data
package pkgregister

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/practicalgo/code/appendix-a/pkgcli/config"
)

type PkgData struct {
	Name     string
	Version  string
	Filename string
	Bytes    io.Reader
}

type PkgRegisterResult struct {
	ID string `json:"id"`
}

func RegisterPackage(ctx context.Context, cliConfig *config.PkgCliConfig, url string, data PkgData) (*PkgRegisterResult, error) {

	p := PkgRegisterResult{}
	payload, contentType, err := createMultiPartMessage(data)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(payload)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/api/packages", reader)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", contentType)
	// TODO this must be documented in the command usage
	authToken := os.Getenv("X_AUTH_TOKEN")
	if len(authToken) != 0 {
		r.Header.Set("X-Auth-Token", authToken)
	}

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return &p, err
	}
	cliConfig.Logger.Debug().Str("pkg_register_response", string(respData))
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respData))
	}
	err = json.Unmarshal(respData, &p)
	return &p, err
}
