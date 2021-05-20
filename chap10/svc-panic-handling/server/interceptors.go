package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type wrappedServerStream struct {
	grpc.ServerStream
}

func (s wrappedServerStream) SendMsg(m interface{}) error {
	log.Printf("Send msg called: %T", m)
	return s.ServerStream.SendMsg(m)
}

func (s wrappedServerStream) RecvMsg(m interface{}) error {
	log.Printf("Waiting to receive a message: %T", m)
	return s.ServerStream.RecvMsg(m)
}

func loggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Print("No metadata")
	}
	log.Printf("Method:%s, Error:%v, Request-Id:%s",
		info.FullMethod,
		err,
		md.Get("Request-Id"),
	)
	return resp, err
}

func loggingStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	serverStream := wrappedServerStream{
		ServerStream: stream,
	}
	err := handler(srv, serverStream)
	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Print("No metadata")
	}
	log.Printf("Method:%s, Error:%v, Request-Id:%s",
		info.FullMethod,
		err,
		md.Get("Request-Id"),
	)
	return err
}

func metricUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	end := time.Now()
	log.Printf("Method:%s, Duration:%s",
		info.FullMethod,
		end.Sub(start),
	)
	return resp, err
}

func metricStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {

	start := time.Now()
	err := handler(srv, stream)
	end := time.Now()
	log.Printf("Method:%s, Duration:%s",
		info.FullMethod,
		end.Sub(start),
	)
	return err
}

func panicUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v\n", r)
			err = status.Error(
				codes.Internal,
				"Unexpected error happened",
			)
		}
	}()
	resp, err = handler(ctx, req)
	return
}

func panicStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v\n", r)
			err = status.Error(
				codes.Internal,
				"Unexpected error happened",
			)
		}
	}()
	serverStream := wrappedServerStream{
		ServerStream: stream,
	}
	err = handler(srv, serverStream)
	return
}
