package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	svc "github.com/practicalgo/code/chap9/bidi-streaming/service"
	"google.golang.org/grpc"
)

func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}

func getUserServiceClient(conn *grpc.ClientConn) svc.UsersClient {
	return svc.NewUsersClient(conn)
}

func setupChat(r io.Reader, w io.Writer, c svc.UsersClient) error {

	stream, err := c.GetHelp(context.Background())
	if err != nil {
		return err
	}
	for {
		scanner := bufio.NewScanner(r)
		prompt := "Request: "
		fmt.Fprint(w, prompt)

		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return err
		}
		msg := scanner.Text()
		if msg == "quit" {
			break
		}
		request := svc.UserHelpRequest{
			Request: msg,
		}
		err := stream.Send(&request)
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("Response: %s\n", resp.Response)
	}
	return stream.CloseSend()
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(
			"Must specify a gRPC server address",
		)
	}
	conn, err := setupGrpcConn(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := getUserServiceClient(conn)
	err = setupChat(os.Stdin, os.Stdout, c)
	if err != nil {
		log.Fatal(err)
	}
}
