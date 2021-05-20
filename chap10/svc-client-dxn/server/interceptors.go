package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func clientDisconnectStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) (err error) {

	ch := make(chan error)

	go func() {
		err = handler(srv, stream)
		ch <- err
	}()

	select {
	case <-stream.Context().Done():
		err = status.Error(
			codes.Canceled,
			fmt.Sprintf("%s: Request canceled", info.FullMethod),
		)
		return
	case <-ch:

	}
	return
}

func clientDisconnectUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	var resp interface{}
	var err error

	ch := make(chan error)

	go func() {
		resp, err = handler(ctx, req)
		ch <- err
	}()

	select {
	case <-ctx.Done():
		err = status.Error(
			codes.Canceled,
			fmt.Sprintf(
				"%s: Request canceled",
				info.FullMethod,
			),
		)
		return resp, err
	case <-ch:

	}
	return resp, err
}
