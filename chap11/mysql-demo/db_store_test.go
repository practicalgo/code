package main

import (
	"database/sql"
	"testing"
)

var testDb *sql.DB

func TestQueryDb(t *testing.T) {

	config := appConfig{
		db: testDb,
	}

	err := updateDb(
		config,
		pkgRow{
			OwnerId:       2,
			Name:          "pkg",
			Version:       "0.3",
			ObjectStoreId: "pkg-0.3-pkg-0.3.tar.gz",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	results, err := queryDb(
		config,
		pkgQueryParams{
			ownerId: 2,
			version: "0.3",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"Expected: 1 row, Got: %d", len(results),
		)
	}

	pkg := results[0]
	t.Logf("%#v", pkg)
	if pkg.RepoURL.Valid {
		t.Fatal("Empty Repo URL expected")
	}

}

func TestQueryDbWithRepo(t *testing.T) {

	config := appConfig{
		db: testDb,
	}

	repoUrl := "https://github.com/practicalgo/code"

	err := updateDb(
		config,
		pkgRow{
			OwnerId:       2,
			Name:          "pkg",
			Version:       "0.4",
			ObjectStoreId: "pkg-0.4-pkg-0.4.tar.gz",
			RepoURL: sql.NullString{
				String: repoUrl,
				Valid:  true,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	results, err := queryDb(
		config,
		pkgQueryParams{
			ownerId: 2,
			version: "0.4",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"Expected: 1 row, Got: %d", len(results),
		)
	}

	pkg := results[0]
	if !pkg.RepoURL.Valid {
		t.Fatal("Repo URL expected")
	}

	if pkg.RepoURL.String != repoUrl {
		t.Fatalf("Expected: %v, Got: %v", repoUrl, pkg.RepoURL)
	}
}
