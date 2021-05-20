package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type wrappedServerStream struct {
	RecvMsgTimeout time.Duration
	grpc.ServerStream
}

func (s wrappedServerStream) SendMsg(m interface{}) error {
	return s.ServerStream.SendMsg(m)
}

func (s wrappedServerStream) RecvMsg(m interface{}) error {
	ch := make(chan error)
	t := time.NewTimer(s.RecvMsgTimeout)
	go func() {
		log.Printf("Waiting to receive a message: %T", m)
		ch <- s.ServerStream.RecvMsg(m)
	}()

	select {
	case <-t.C:
		return status.Error(
			codes.DeadlineExceeded,
			"Deadline exceeded",
		)
	case err := <-ch:
		return err
	}
}

func timeoutStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	serverStream := wrappedServerStream{
		RecvMsgTimeout: 500 * time.Millisecond,
		ServerStream:   stream,
	}
	err := handler(srv, serverStream)
	return err
}

func timeoutUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	var resp interface{}
	var err error

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	ch := make(chan error)

	go func() {
		resp, err = handler(ctxWithTimeout, req)
		ch <- err
	}()

	select {
	case <-ctxWithTimeout.Done():
		cancel()
		err = status.Error(
			codes.DeadlineExceeded,
			fmt.Sprintf("%s: Deadline exceeded", info.FullMethod),
		)
		return resp, err
	case <-ch:

	}
	return resp, err
}
