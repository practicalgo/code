package main

import (
	"database/sql"
	"time"
)

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
	RepoURL       sql.NullString
	Created       time.Time
}

type appConfig struct {
	db *sql.DB
}
