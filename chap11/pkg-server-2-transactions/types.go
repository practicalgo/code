package main

type pkgRegisterResponse struct {
	ID string `json:"id"`
}

type pkgQueryResponse struct {
	ID string `json:"id"`
}

type pkgQueryParams struct {
	name    string
	version string
	ownerId int
}

type pkgRow struct {
	OwnerId       int
	Name          string
	Version       string
	ObjectStoreId string
	Created       string
}
