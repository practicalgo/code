package interceptors

import (
	"context"
	"time"

	"github.com/practicalgo/code/appendix-a/grpc-server/server/config"
	"github.com/practicalgo/code/appendix-a/grpc-server/server/telemetry"
	"google.golang.org/grpc"
)

func MetricUnaryInterceptor(config *config.AppConfig) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		metricName := "userssvc.grpc.unary_latency"
		duration := time.Since(startTime).Seconds()
		config.Metrics.ReportLatency(
			metricName,
			telemetry.DurationMetric{
				Method:     info.FullMethod,
				DurationMs: duration,
				Success:    err == nil,
			},
		)
		return resp, err
	}
}

func MetricStreamInterceptor(config *config.AppConfig) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {

		startTime := time.Now()
		err := handler(srv, stream)
		metricName := "userssvc.grpc.stream_latency"
		duration := time.Since(startTime).Seconds()
		config.Metrics.ReportLatency(
			metricName,
			telemetry.DurationMetric{
				Method:     info.FullMethod,
				DurationMs: duration,
				Success:    err == nil,
			},
		)
		return err
	}
}
