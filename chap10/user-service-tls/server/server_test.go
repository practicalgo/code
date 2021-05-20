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
	registerServices(s)
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
	creds, err := credentials.NewClientTLSFromFile(
		tlsCertFile,
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	credsOption := grpc.WithTransportCredentials(creds)
	client, err := grpc.DialContext(
		context.Background(),
		"localhost",
		credsOption,
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}
	usersClient := users.NewUsersClient(client)
	resp, err := usersClient.GetUser(
		context.Background(),
		&users.UserGetRequest{
			Email: "jane@doe.com",
			Id:    "foo-bar",
		},
	)

	if err != nil {
		t.Fatal(err)
	}
	if resp.User.FirstName != "jane" {
		t.Errorf(
			"Expected FirstName to be: jane, Got: %s",
			resp.User.FirstName,
		)
	}

}
