package main

import (
	"fmt"
	"net/url"
	"os"

	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
)

func getTestBucket(tmpDir string) (*blob.Bucket, error) {
	myDir, err := os.MkdirTemp(tmpDir, "test-bucket")
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(fmt.Sprintf("file:///%s", myDir))
	if err != nil {
		return nil, err
	}
	opts := fileblob.Options{
		URLSigner: fileblob.NewURLSignerHMAC(
			u,
			[]byte("super secret"),
		),
	}
	return fileblob.OpenBucket(myDir, &opts)
}
