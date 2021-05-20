package main

import (
	"context"
	"errors"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthsvc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func getHealthSvcClient(l *bufconn.Listener) (healthsvc.HealthClient, error) {

	bufconnDialer := func(
		ctx context.Context, addr string,
	) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(
		context.Background(),
		"", grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		return nil, err
	}
	return healthsvc.NewHealthClient(client), nil
}

func TestHealthService(t *testing.T) {

	l := startTestGrpcServer()
	healthClient, err := getHealthSvcClient(l)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := healthClient.Check(
		context.Background(),
		&healthsvc.HealthCheckRequest{},
	)
	if err != nil {
		t.Fatal(err)
	}
	serviceHealthStatus := resp.Status.String()
	if serviceHealthStatus != "SERVING" {
		t.Fatalf(
			"Expected health: SERVING, Got: %s",
			serviceHealthStatus,
		)
	}
}

func TestHealthServiceUsers(t *testing.T) {

	l := startTestGrpcServer()
	healthClient, err := getHealthSvcClient(l)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := healthClient.Check(
		context.Background(),
		&healthsvc.HealthCheckRequest{
			Service: "Users",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	serviceHealthStatus := resp.Status.String()
	if serviceHealthStatus != "SERVING" {
		t.Fatalf(
			"Expected health: SERVING, Got: %s",
			serviceHealthStatus,
		)
	}
}

func TestHealthServiceUnknown(t *testing.T) {

	l := startTestGrpcServer()
	healthClient, err := getHealthSvcClient(l)
	if err != nil {
		t.Fatal(err)
	}

	_, err = healthClient.Check(
		context.Background(),
		&healthsvc.HealthCheckRequest{
			Service: "Repo",
		},
	)
	if err == nil {
		t.Fatalf("Expected non-nil error, Got nil error")
	}
	expectedError := status.Errorf(
		codes.NotFound, "unknown service",
	)
	if !errors.Is(err, expectedError) {
		t.Fatalf(
			"Expected error %v, Got; %v",
			err,
			expectedError,
		)
	}
}

func TestHealthServiceWatch(t *testing.T) {

	l := startTestGrpcServer()
	healthClient, err := getHealthSvcClient(l)
	if err != nil {
		t.Fatal(err)
	}

	client, err := healthClient.Watch(
		context.Background(),
		&healthsvc.HealthCheckRequest{
			Service: "Users",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Recv()
	if err != nil {
		t.Fatalf("Error in Watch: %#v\n", err)
	}
	if resp.Status != healthsvc.HealthCheckResponse_SERVING {
		t.Errorf("Expected SERVING, Got: %#v", resp.Status.String())
	}

	updateServiceHealth(
		h,
		"Users",
		healthsvc.HealthCheckResponse_NOT_SERVING,
	)

	resp, err = client.Recv()
	if err != nil {
		t.Fatalf("Error in Watch: %#v\n", err)
	}
	if resp.Status != healthsvc.HealthCheckResponse_NOT_SERVING {
		t.Errorf("Expected NOT_SERVING, Got: %#v", resp.Status.String())
	}
}
