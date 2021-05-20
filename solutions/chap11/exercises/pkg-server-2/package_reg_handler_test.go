package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "gocloud.dev/blob/fileblob"
)

type pkgData struct {
	Name     string
	Version  string
	Filename string
	Bytes    io.Reader
}

func createMultiPartMessage(data pkgData) ([]byte, string, error) {
	var b bytes.Buffer
	var err error
	var fw io.Writer

	mw := multipart.NewWriter(&b)

	fw, err = mw.CreateFormField("name")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Name)

	fw, err = mw.CreateFormField("version")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Version)

	fw, err = mw.CreateFormFile("filedata", data.Filename)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(fw, data.Bytes)
	if err != nil {
		return nil, "", err
	}
	err = mw.Close()
	if err != nil {
		return nil, "", err
	}

	contentType := mw.FormDataContentType()
	return b.Bytes(), contentType, nil
}

func TestPackageRegHandler(t *testing.T) {
	packageBucket, err := getTestBucket(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer packageBucket.Close()

	testC, testDb, err := getTestDb()
	if err != nil {
		t.Fatal(err)
	}
	defer testC.Terminate(context.Background())

	config := appConfig{
		logger: log.New(
			os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile,
		),
		packageBucket: packageBucket,
		db:            testDb,
	}

	mux := http.NewServeMux()
	setupHandlers(mux, config)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	p := pkgData{
		Name:     "mypackage",
		Version:  "0.1",
		Filename: "mypackage-0.1.tar.gz",
		Bytes:    strings.NewReader("data"),
	}

	payload, contentType, err := createMultiPartMessage(p)
	if err != nil {
		t.Fatal(err)
	}
	reader := bytes.NewReader(payload)
	r, err := http.Post(ts.URL+"/api/packages", contentType, reader)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	respData, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	resp := pkgRegisterResponse{}
	err = json.Unmarshal(respData, &resp)
	if err != nil {
		t.Fatal(err)
	}
	expectedPackageId := "mypackage-0.1-mypackage-0.1.tar.gz"
	if resp.ID != expectedPackageId {
		t.Fatalf(
			"Expected version to be %s, Got: %s\n",
			expectedPackageId,
			resp.ID,
		)
	}
}
