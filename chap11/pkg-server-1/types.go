package main

import "io"

type pkgData struct {
	Name     string
	Version  string
	Filename string
	Bytes    io.Reader
}

type pkgRegisterResponse struct {
	ID string `json:"id"`
}
