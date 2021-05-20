package main

import (
	"context"
	"log"
	"net"
	"testing"

	users "github.com/practicalgo/code/chap10/client-resiliency/service"
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

func startTestGrpcServer() *bufconn.Listener {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	users.RegisterUsersServer(s, &dummyUserService{})
	go func() {
		log.Fatal(s.Serve(l))
	}()
	return l
}

func TestGetUser(t *testing.T) {

	l := startTestGrpcServer()

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
	result, _, err := getUser(
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
