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

	svc "github.com/practicalgo/code/chap9/interceptor-chain/service"
	"google.golang.org/grpc"
)

type userService struct {
	svc.UnimplementedUsersServer
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
	for {

		request, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("Request receieved: %s\n", request.Request)
		response := svc.UserHelpReply{
			Response: request.Request,
		}
		err = stream.Send(&response)
		if err != nil {
			return err
		}
	}
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
		grpc.ChainUnaryInterceptor(
			metricUnaryInterceptor,
			loggingUnaryInterceptor,
		),
		grpc.ChainStreamInterceptor(
			metricStreamInterceptor,
			loggingStreamInterceptor,
		),
	)
	registerServices(s)
	log.Fatal(startServer(s, lis))
}
