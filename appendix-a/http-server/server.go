package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"

	_ "github.com/go-sql-driver/mysql"

	users "github.com/practicalgo/code/appendix-a/grpc-server/service"
	"github.com/practicalgo/code/appendix-a/http-server/config"
	"github.com/practicalgo/code/appendix-a/http-server/handlers"
	"github.com/practicalgo/code/appendix-a/http-server/middleware"
	"github.com/practicalgo/code/appendix-a/http-server/storage"
	"github.com/practicalgo/code/appendix-a/http-server/telemetry"
	"google.golang.org/grpc"
)

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

func getUserServiceClient(conn *grpc.ClientConn) users.UsersClient {
	return users.NewUsersClient(conn)
}

func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
}

func main() {

	version := "0.1"

	// Object storage
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

	// Database
	dbAddr := os.Getenv("DB_ADDR")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	if len(dbAddr) == 0 || len(dbName) == 0 || len(dbUser) == 0 || len(dbPassword) == 0 {
		log.Fatal(
			"Must specfy DB details - DB_ADDR, DB_NAME, DB_USER, DB_PASSWORD",
		)
	}

	db, err := storage.GetDatabaseConn(
		dbAddr, dbName,
		dbUser, dbPassword,
	)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Telemetry stuff
	logLevel := os.Getenv("LOG_LEVEL")
	if len(logLevel) == 0 {
		logLevel = "1"
	}
	logLevelInt, err := strconv.Atoi(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	statsdAddr := os.Getenv("STATSD_ADDR")
	if len(statsdAddr) == 0 {
		log.Fatal("Specify STATSD_ADDR")
	}

	metricReporter, err := telemetry.InitMetrics(statsdAddr)
	if err != nil {
		log.Fatal(err)
	}

	jaegerAddr := os.Getenv("JAEGER_ADDR")
	if len(statsdAddr) == 0 {
		log.Fatal("Specify JAEGER_ADDR")
	}

	err = telemetry.InitTracing(jaegerAddr)
	if err != nil {
		log.Fatal(err)
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	usersSvcAddr := os.Getenv("USERS_SVC_ADDR")
	if len(usersSvcAddr) == 0 {
		log.Fatal("Specify USERS_SVC_ADDR")
	}
	grpcConn, err := setupGrpcConn(usersSvcAddr)
	if err != nil {
		log.Fatal(err)
	}
	grpcUsers := getUserServiceClient(grpcConn)

	config := config.AppConfig{
		Logger:        telemetry.InitLogging(os.Stdout, version, logLevelInt),
		Metrics:       metricReporter,
		PackageBucket: packageBucket,
		Db:            db,
		UsersSvc:      grpcUsers,
	}

	mux := http.NewServeMux()
	handlers.SetupHandlers(mux, &config)
	m := middleware.LoggingMiddleware(
		&config,
		middleware.MetricMiddleware(
			&config,
			middleware.TracingMiddleware(
				&config, mux,
			),
		),
	)
	config.Logger.Info().Msg("Starting HTTP server")
	log.Fatal(http.ListenAndServe(listenAddr, m))
}
