package types

import (
	"net/http"

	"github.com/practicalgo/code/appendix-a/http-server/config"
)

type PkgRegisterResponse struct {
	ID string `json:"id"`
}

type PkgQueryResponse struct {
	ID string `json:"id"`
}

type PkgQueryParams struct {
	Name    string
	Version string
	OwnerId int
}

type PkgRow struct {
	OwnerId       int
	Name          string
	Version       string
	ObjectStoreId string
	Created       string
}

type App struct {
	Config  *config.AppConfig
	Handler func(w http.ResponseWriter, r *http.Request, config *config.AppConfig)
}

func (a App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Handler(w, r, a.Config)
}
