package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "gocloud.dev/blob/fileblob"
)

func TestPackageGetHandler(t *testing.T) {
	packageBucket, err := getTestBucket(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer packageBucket.Close()

	// create a test object
	err = packageBucket.WriteAll(
		context.Background(),
		"test-object-id",
		[]byte("test-data"),
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	config := appConfig{
		logger: log.New(
			os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile,
		),
		packageBucket: packageBucket,
	}

	mux := http.NewServeMux()
	setupHandlers(mux, config)

	ts := httptest.NewServer(mux)
	defer ts.Close()
	var redirectUrl string
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirectUrl = req.URL.String()
			return errors.New("no redirect")
		},
	}

	_, err = client.Get(ts.URL + "/api/packages?id=test-object-id")
	if err == nil {
		t.Fatal("Expected error: no redirect, Got nil")
	}
	if !strings.HasPrefix(redirectUrl, "file:///") {
		t.Fatalf("Expected redirect url to start with file:///, got: %v", redirectUrl)
	}
}
