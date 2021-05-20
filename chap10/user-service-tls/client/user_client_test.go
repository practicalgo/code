package main

import (
	"context"
	"log"
	"net"
	"testing"

	users "github.com/practicalgo/code/chap10/user-service-tls/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/test/bufconn"
)

var (
	tlsCertFile string = "../tls/server.crt"
	tlsKeyFile  string = "../tls/server.key"
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
	creds, err := credentials.NewServerTLSFromFile(
		tlsCertFile,
		tlsKeyFile,
	)
	if err != nil {
		log.Fatal(err)
	}
	credsOption := grpc.Creds(creds)
	s := grpc.NewServer(credsOption)
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
	creds, err := credentials.NewClientTLSFromFile(
		tlsCertFile,
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	credsOption := grpc.WithTransportCredentials(creds)
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost",
		credsOption,
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
