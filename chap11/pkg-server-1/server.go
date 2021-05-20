package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"
)

type appConfig struct {
	logger        *log.Logger
	packageBucket *blob.Bucket
}

type app struct {
	config  appConfig
	handler func(w http.ResponseWriter, r *http.Request, config appConfig)
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler(w, r, a.config)
}

func setupHandlers(mux *http.ServeMux, config appConfig) {
	mux.Handle(
		"/api/packages",
		&app{config: config, handler: packageHandler},
	)
}

func getBucket(
	bucketName, s3Address, s3Region string,
) (*blob.Bucket, error) {

	urlString := fmt.Sprintf("s3://%s?", bucketName)
	if len(s3Region) != 0 {
		urlString += fmt.Sprintf("region=%s&", s3Region)
	}

	if len(s3Address) != 0 {
		urlString += fmt.Sprintf("endpoint=%s&"+
			"disableSSL=true&"+
			"s3ForcePathStyle=true",
			s3Address,
		)
	}
	return blob.OpenBucket(context.Background(), urlString)
}

func main() {

	bucketName := os.Getenv("BUCKET_NAME")
	if len(bucketName) == 0 {
		log.Fatal("Specify Object Storage bucket - BUCKET_NAME")
	}
	s3Address := os.Getenv("S3_ADDR")
	awsRegion := os.Getenv("AWS_DEFAULT_REGION")

	if len(s3Address) == 0 && len(awsRegion) == 0 {
		log.Fatal(
			"Assuming AWS S3 service. Specify AWS_DEFAULT_REGION",
		)
	}

	packageBucket, err := getBucket(
		bucketName, s3Address, awsRegion,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer packageBucket.Close()

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	config := appConfig{
		logger: log.New(
			os.Stdout, "",
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		packageBucket: packageBucket,
	}

	mux := http.NewServeMux()
	setupHandlers(mux, config)

	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
