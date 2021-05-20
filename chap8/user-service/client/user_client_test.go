package main

import (
	"context"
	"log"
	"net"
	"testing"

	users "github.com/practicalgo/code/chap8/user-service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type dummyUserService struct {
	users.UnimplementedUsersServer
}

func (s *dummyUserService) GetUser(
	ctx context.Context,
	in *users.UserGetRequest,
) (*users.UserGetReply, error) {
	u := users.User{
		Id:        "user-123-a",
		FirstName: "jane",
		LastName:  "doe",
		Age:       36,
	}
	return &users.UserGetReply{User: &u}, nil
}

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	users.RegisterUsersServer(s, &dummyUserService{})
	go func() {
		err := s.Serve(l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}

func TestGetUser(t *testing.T) {

	s, l := startTestGrpcServer()
	defer s.GracefulStop()

	bufconnDialer := func(
		ctx context.Context, addr string,
	) (net.Conn, error) {
		return l.Dial()
	}

	conn, err := grpc.DialContext(
		context.Background(),
		"", grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}

	c := getUserServiceClient(conn)
	result, err := getUser(
		c,
		&users.UserGetRequest{Email: "jane@doe.com"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if result.User.FirstName != "jane" ||
		result.User.LastName != "doe" {
		t.Fatalf(
			"Expected: jane doe, Got: %s %s",
			result.User.FirstName,
			result.User.LastName,
		)
	}
}
