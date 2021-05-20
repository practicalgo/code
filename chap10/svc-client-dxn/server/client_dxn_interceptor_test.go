package main

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	svc "github.com/practicalgo/code/chap10/svc-client-dxn/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestClientDxnInterceptor(t *testing.T) {
	req := svc.UserGetRequest{}
	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "Users.GetUser",
	}
	testUnaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(200 * time.Millisecond)
		return svc.UserGetReply{}, nil
	}

	incomingContext, cancel := context.WithTimeout(
		context.Background(),
		100*time.Millisecond,
	)
	defer cancel()

	_, err := clientDisconnectUnaryInterceptor(
		incomingContext,
		&req,
		unaryInfo,
		testUnaryHandler,
	)
	expectedErr := status.Errorf(
		codes.Canceled,
		"Users.GetUser: Request canceled",
	)
	if !errors.Is(err, expectedErr) {
		t.Errorf(
			"Expected error: %v Got: %v\n",
			expectedErr,
			err,
		)
	}
}

type testStream struct {
	CancelFunc context.CancelFunc
	grpc.ServerStream
}

func (s testStream) Context() context.Context {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		100*time.Millisecond,
	)
	s.CancelFunc = cancel
	return ctx
}

func TestStreamingClientDxnInterceptor(t *testing.T) {

	streamInfo := &grpc.StreamServerInfo{
		FullMethod:     "Users.GetUser",
		IsClientStream: true,
		IsServerStream: true,
	}

	testStream := testStream{}

	testHandler := func(srv interface{}, stream grpc.ServerStream) (err error) {
		time.Sleep(200 * time.Millisecond)
		for {
			m := svc.UserHelpRequest{}
			err := stream.RecvMsg(&m)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err

			}
			r := svc.UserHelpReply{}
			err = stream.SendMsg(&r)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err

			}
		}
		return nil
	}

	err := clientDisconnectStreamInterceptor(
		"test",
		testStream,
		streamInfo,
		testHandler,
	)
	expectedErr := status.Errorf(
		codes.Canceled,
		"Users.GetUser: Request canceled",
	)
	if !errors.Is(err, expectedErr) {
		t.Errorf(
			"Expected error: %v Got: %v\n",
			expectedErr,
			err,
		)
	}
}
