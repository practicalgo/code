package main

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"
)

var testDb *sql.DB

func TestQueryDb(t *testing.T) {

	config := appConfig{
		logger: log.New(
			os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile,
		),
		db: testDb,
	}

	// update package metadata for the test object
	err := updateDb(
		config,
		pkgRow{
			OwnerId:       1,
			Name:          "pkg",
			Version:       "0.2",
			ObjectStoreId: "pkg-0.2-pkg-0.2.tar.gz",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// update package metadata for the test object
	err = updateDb(
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

	layout := "2006-01-02 15:04:05"
	created := results[0].Created
	parsedTime, err := time.Parse(layout, created)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", parsedTime.Local().String())
}
