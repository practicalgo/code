package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/practicalgo/code/appendix-b/grpc-server/server/config"
	"github.com/practicalgo/code/appendix-b/grpc-server/server/interceptors"
	"github.com/practicalgo/code/appendix-b/grpc-server/server/telemetry"
	users "github.com/practicalgo/code/appendix-b/grpc-server/service"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func registerServices(s *grpc.Server, config *config.AppConfig) {
	u := userService{}
	users.RegisterUsersServer(s, &u)
}

func startServer(s *grpc.Server, l net.Listener) error {
	return s.Serve(l)
}

func main() {

	version := "0.1"

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":50051"
	}

	// Telemetry stuff
	jaegerAddr := os.Getenv("JAEGER_ADDR")
	if len(jaegerAddr) == 0 {
		log.Fatal("Specify JAEGER_ADDR")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if len(logLevel) == 0 {
		logLevel = "1"
	}
	logLevelInt, err := strconv.Atoi(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	err = telemetry.InitTracing(jaegerAddr)
	if err != nil {
		log.Fatal("Error initializing tracing", err)
	}
	statsdAddr := os.Getenv("STATSD_ADDR")
	if len(statsdAddr) == 0 {
		log.Fatal("Specify STATSD_ADDR")
	}
	metricReporter, err := telemetry.InitMetrics(statsdAddr)
	if err != nil {
		log.Fatal(err)
	}
	config := config.AppConfig{
		Logger:  telemetry.InitLogging(os.Stdout, version, logLevelInt),
		Metrics: metricReporter,
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.MetricUnaryInterceptor(&config),
			interceptors.LoggingUnaryInterceptor(&config),
			otelgrpc.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			interceptors.MetricStreamInterceptor(&config),
			interceptors.LoggingStreamInterceptor(&config),
			otelgrpc.StreamServerInterceptor(),
		),
	)

	registerServices(s, &config)
	config.Logger.Info().Msg("Starting gRPC server")
	log.Fatal(startServer(s, lis))
}
