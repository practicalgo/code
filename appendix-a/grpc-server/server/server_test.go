package main

import (
	"context"
	"log"
	"net"
	"testing"

	users "github.com/practicalgo/code/appendix-a/grpc-server/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	registerServices(s, nil)
	go func() {
		err := startServer(s, l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}
func TestUserService(t *testing.T) {

	s, l := startTestGrpcServer()
	defer s.GracefulStop()

	bufconnDialer := func(
		ctx context.Context, addr string,
	) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(
		context.Background(),
		"", grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}
	usersClient := users.NewUsersClient(client)
	resp, err := usersClient.GetUser(
		context.Background(),
		&users.UserGetRequest{
			Auth: "token-123",
		},
	)

	if err != nil {
		t.Fatal(err)
	}
	if resp.User.Id < 1 || resp.User.Id > 5 {
		t.Errorf(
			"Expected Id to be in the range [1, 5]. Got: %d", resp.User.Id,
		)
	}
}
