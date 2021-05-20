package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"

	_ "github.com/go-sql-driver/mysql"
)

type appConfig struct {
	logger        *log.Logger
	packageBucket *blob.Bucket
	db            *sql.DB
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
	bucketName,
	s3Address,
	s3Region string,
) (*blob.Bucket, error) {

	urlString := fmt.Sprintf("s3://%s?", bucketName)
	if len(s3Region) == 0 {
		s3Region = "local"
	}
	urlString += fmt.Sprintf("region=%s&", s3Region)

	if len(s3Address) != 0 {
		urlString += fmt.Sprintf(
			"endpoint=%s&"+
				"disableSSL=true&"+
				"s3ForcePathStyle=true",
			s3Address,
		)
	}
	return blob.OpenBucket(
		context.Background(),
		urlString,
	)
}

func main() {

	bucketName := os.Getenv("BUCKET_NAME")
	if len(bucketName) == 0 {
		log.Fatal("Specify BUCKET_NAME.")
	}
	s3Address := os.Getenv("S3_ADDR")
	s3Region := os.Getenv("S3_REGION")

	packageBucket, err := getBucket(
		bucketName, s3Address, s3Region,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer packageBucket.Close()

	dbAddr := os.Getenv("DB_ADDR")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	if len(dbAddr) == 0 || len(dbName) == 0 || len(dbUser) == 0 || len(dbPassword) == 0 {
		log.Fatal(
			"Must specfy DB details - DB_ADDR, DB_NAME, DB_USER, DB_PASSWORD",
		)
	}

	db, err := getDatabaseConn(
		dbAddr, dbName,
		dbUser, dbPassword,
	)

	if err != nil {
		log.Fatal(err)
	}

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
		db:            db,
	}

	mux := http.NewServeMux()
	setupHandlers(mux, config)

	log.Fatal(http.ListenAndServe(listenAddr, mux))
	db.Close()
}
