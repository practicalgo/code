package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"strings"

	users "github.com/practicalgo/code/chap8/user-service/service-v2"
	"google.golang.org/grpc"
)

type userService struct {
	users.UnimplementedUsersServer
}

func (s *userService) GetUser(
	ctx context.Context,
	in *users.UserGetRequest,
) (*users.UserGetReply, error) {
	log.Printf("Received request for user with Email: %s Id: %s\n", in.Email, in.Id)
	u := users.User{}

	if len(in.Email) != 0 {
		components := strings.Split(in.Email, "@")
		if len(components) != 2 {
			return nil, errors.New("invalid email address")
		}
		u.FirstName = components[0]
		u.LastName = components[1]
	}
	u.Age = 36
	if len(in.Id) != 0 {
		u.Id = in.Id
	}
        u.Location = "AU"
	return &users.UserGetReply{User: &u}, nil
}

func registerServices(s *grpc.Server) {
	users.RegisterUsersServer(s, &userService{})
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
	s := grpc.NewServer()
	registerServices(s)
	log.Fatal(startServer(s, lis))
}
