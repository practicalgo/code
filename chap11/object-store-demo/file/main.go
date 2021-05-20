package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
)

func main() {
	myDir, err := os.MkdirTemp("", "test-bucket")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(myDir)

	ctx := context.Background()
	bucket, err := blob.OpenBucket(
		ctx,
		"file:///"+myDir,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer bucket.Close()

	err = bucket.WriteAll(
		ctx,
		"object-id",
		[]byte("Hello world"),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	data, err := bucket.ReadAll(
		ctx,
		"object-id",
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data: %s\n", string(data))
}
