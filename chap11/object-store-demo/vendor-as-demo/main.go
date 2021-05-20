package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/s3"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"
)

func main() {

	bucketName := "practicalgo-echorand"
	testBucket, err := blob.OpenBucket(
		context.Background(),
		fmt.Sprintf("s3://%s", bucketName),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer testBucket.Close()

	var s3Svc *s3.S3
	if !testBucket.As(&s3Svc) {
		log.Fatal("Couldn't convert type to underlying bucket type")
	}
	_, err = s3Svc.HeadBucket(
		&s3.HeadBucketInput{
			Bucket: &bucketName,
		},
	)
	if err != nil {
		log.Fatalf(
			"Bucket doesn't exist, or insufficient permissions: %v\n",
			err,
		)
	}
}
