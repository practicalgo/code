package main

import (
	"context"
	"fmt"
	"log"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"
)

func main() {

	ctx := context.Background()
	bucket, err := blob.OpenBucket(
		ctx,
		"s3://test-bucket?"+
			"endpoint=127.0.0.1:9000&"+
			"region=local&"+
			"disableSSL=true&"+
			"s3ForcePathStyle=true",
	)

	if err != nil {
		log.Fatal(err)
	}
	defer bucket.Close()

	err = bucket.WriteAll(
		ctx,
		"object-id-1",
		[]byte("Hello world"),
		nil,
	)
	if err != nil {
		log.Fatalf("Error creating object: %v", err)
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
