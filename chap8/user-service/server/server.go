package main

import (
	"context"
	"log"
	"net"
	"os"

	users "github.com/practicalgo/code/chap8/user-service/service"
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
	u := users.User{Id: "user-123-a", FirstName: "Jane", LastName: "Doe", Age: 36}
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
