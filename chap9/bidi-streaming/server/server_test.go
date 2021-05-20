package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	svc "github.com/practicalgo/code/chap9/bidi-streaming/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func startTestGrpcServer() *bufconn.Listener {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	registerServices(s)
	go func() {
		log.Fatal(startServer(s, l))
	}()
	return l
}
func TestGetHelp(t *testing.T) {

	l := startTestGrpcServer()

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
	usersClient := svc.NewUsersClient(client)
	stream, err := usersClient.GetHelp(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i <= 5; i++ {
		msg := fmt.Sprintf("Hello: %d", i)
		request := svc.UserHelpRequest{
			Request: msg,
		}
		err := stream.Send(&request)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := stream.Recv()
		if err != nil {
			t.Fatal(err)
		}
		if resp.Response != msg {
			t.Errorf(
				"Expected Response to be: %s, Got: %s",
				msg,
				resp.Response,
			)
		}
	}

}
