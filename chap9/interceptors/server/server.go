package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	svc "github.com/practicalgo/code/chap9/interceptors/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type userService struct {
	svc.UnimplementedUsersServer
}

func logMessage(
	ctx context.Context,
	method string,
	latency time.Duration,
	err error,
) {
	var requestId string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Print("No metadata")
	} else {
		if len(md.Get("Request-Id")) != 0 {
			requestId = md.Get("Request-Id")[0]
		}
	}
	log.Printf("Method:%s, Duration:%s, Error:%v, Request-Id:%s",
		method,
		latency,
		err,
		requestId,
	)
}

func loggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	logMessage(ctx, info.FullMethod, time.Since(start), err)
	return resp, err
}

func loggingStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {

	start := time.Now()
	err := handler(srv, stream)
	ctx := stream.Context()
	logMessage(ctx, info.FullMethod, time.Since(start), err)
	return err
}

func (s *userService) GetUser(
	ctx context.Context,
	in *svc.UserGetRequest,
) (*svc.UserGetReply, error) {

	log.Printf(
		"Received request for user with Email: %s Id: %s\n",
		in.Email,
		in.Id,
	)
	components := strings.Split(in.Email, "@")
	if len(components) != 2 {
		return nil, errors.New("invalid email address")
	}
	u := svc.User{
		Id:        in.Id,
		FirstName: components[0],
		LastName:  components[1],
		Age:       36,
	}
	return &svc.UserGetReply{User: &u}, nil
}

func (s *userService) GetHelp(
	stream svc.Users_GetHelpServer,
) error {
	log.Println("Client connected")
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("Request received: %s\n", request.Request)
		response := svc.UserHelpReply{
			Response: request.Request,
		}
		err = stream.Send(&response)
		if err != nil {
			return err
		}
	}
	log.Println("Client disconnected")
	return nil
}

func registerServices(s *grpc.Server) {
	svc.RegisterUsersServer(s, &userService{})
}

func startServer(s *grpc.Server, l net.Listener) error {
	return s.Serve(l)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":50051"
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(loggingUnaryInterceptor),
		grpc.StreamInterceptor(loggingStreamInterceptor),
	)
	registerServices(s)
	log.Fatal(startServer(s, lis))
}
