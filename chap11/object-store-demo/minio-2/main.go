package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"gocloud.dev/blob/s3blob"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewEnvCredentials(),
		Endpoint:         aws.String("http://127.0.0.1:9000"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	bucket, err := s3blob.OpenBucket(
		ctx,
		sess,
		"test-bucket",
		nil,
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
