package main

import (
	"log"

	svc "github.com/practicalgo/code/chap10/server-healthcheck/service"
	"google.golang.org/grpc"
	healthz "google.golang.org/grpc/health"
	healthsvc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"
)

var h *healthz.Server

func startTestGrpcServer() *bufconn.Listener {
	h = healthz.NewServer()
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	registerServices(s, h)
	updateServiceHealth(
		h,
		svc.Users_ServiceDesc.ServiceName,
		healthsvc.HealthCheckResponse_SERVING,
	)
	go func() {
		log.Fatal(startServer(s, l))
	}()
	return l
}
