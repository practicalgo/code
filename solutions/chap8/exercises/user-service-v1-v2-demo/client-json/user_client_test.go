package main

import (
	"context"
	"log"
	"net"
	"strings"
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
	components := strings.Split(in.Email, "@")
	u := users.User{
		Id:        in.Id,
		FirstName: components[0],
		LastName:  components[1],
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

	u, err := createUserRequest(`{"email":"john@doe.com","id":"user-123"}`)
	if err != nil {
		t.Fatal(err)
	}
	c := getUserServiceClient(conn)
	result, err := getUser(
		c,
		u,
	)
	if err != nil {
		t.Fatal(err)
	}

	respData, err := getUserResponseJson(result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"id":"user-123", "firstName":"john", "lastName":"doe.com", "age":36}`
	if !strings.Contains(string(respData), expected) {
		t.Fatalf("Expected: %s to contain :%s\n", string(respData), expected)
	}
}
