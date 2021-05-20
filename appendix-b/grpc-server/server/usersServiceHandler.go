package main

import (
	"context"
	"math/rand"

	"github.com/practicalgo/code/appendix-b/grpc-server/server/config"
	users "github.com/practicalgo/code/appendix-b/grpc-server/service"
)

type userService struct {
	users.UnimplementedUsersServer
	config config.AppConfig
}

func (s *userService) GetUser(
	ctx context.Context,
	in *users.UserGetRequest,
) (*users.UserGetReply, error) {
	s.config.Logger.Printf(
		"Received request for user verification: %s\n",
		in.Auth,
	)
	u := users.User{
		Id: rand.Int31n(4) + 1,
	}
	return &users.UserGetReply{User: &u}, nil
}
