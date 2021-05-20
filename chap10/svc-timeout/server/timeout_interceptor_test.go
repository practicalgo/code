package main

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"
	"time"

	svc "github.com/practicalgo/code/chap10/svc-timeout/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryTimeOutInterceptor(t *testing.T) {
	req := svc.UserGetRequest{}
	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "Users.GetUser",
	}
	testUnaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(500 * time.Millisecond)
		return svc.UserGetReply{}, nil
	}

	_, err := timeoutUnaryInterceptor(
		context.Background(),
		&req,
		unaryInfo,
		testUnaryHandler,
	)
	if err == nil {
		t.Fatal(err)
	}
	expectedErr := status.Errorf(
		codes.DeadlineExceeded,
		"Users.GetUser: Deadline exceeded",
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
	grpc.ServerStream
}

func (s testStream) SendMsg(m interface{}) error {
	log.Println("Test Stream - SendMsg")
	return nil
}

func (s testStream) RecvMsg(m interface{}) error {
	log.Println("Test Stream - RecvMsg - Going to sleep")
	time.Sleep(700 * time.Millisecond)
	return nil
}

func TestStreamingTimeOutInterceptor(t *testing.T) {

	streamInfo := &grpc.StreamServerInfo{
		FullMethod:     "Users.GetUser",
		IsClientStream: true,
		IsServerStream: true,
	}

	testStream := testStream{}

	testHandler := func(srv interface{}, stream grpc.ServerStream) (err error) {
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

	err := timeoutStreamInterceptor(
		"test",
		testStream,
		streamInfo,
		testHandler,
	)
	expectedErr := status.Errorf(
		codes.DeadlineExceeded,
		"Deadline exceeded",
	)
	if !errors.Is(err, expectedErr) {
		t.Errorf(
			"Expected error: %v Got: %v\n",
			expectedErr,
			err,
		)
	}
}
