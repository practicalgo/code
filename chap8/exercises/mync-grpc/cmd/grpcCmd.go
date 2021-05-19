package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	users "github.com/practicalgo/code/chap8/user-service/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type grpcConfig struct {
	server string
	method string
	body   string
}

func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}

func getUserServiceClient(conn *grpc.ClientConn) users.UsersClient {
	return users.NewUsersClient(conn)
}

func getUser(
	client users.UsersClient,
	u *users.UserGetRequest,
) (*users.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func handleRequest(
	addr string,
	method string,
	body string,
) error {
	client, err := setupGrpcConn(addr)
	if err != nil {
		return err
	}
	data := users.UserGetRequest{}
	input := []byte(body)
	err = protojson.Unmarshal(input, &data)
	if err != nil {
		log.Fatal(err)
	}

	u := getUserServiceClient(client)
	resp, err := getUser(u, &data)
	if err != nil {
		return err
	}
	log.Printf("%#v\n", resp)
	return nil
}

func HandleGrpc(w io.Writer, args []string) error {
	log.Printf("%#v\n", args)
	c := grpcConfig{}
	fs := flag.NewFlagSet("grpc", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.method, "method", "", "Method to call")
	fs.StringVar(&c.body, "body", "", "Body of request")
	fs.Usage = func() {
		var usageString = `
grpc: A gRPC client.

grpc: <options> server`
		fmt.Fprint(w, usageString)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}
	c.server = fs.Arg(0)
	fmt.Fprintln(w, "Executing grpc command")
	handleRequest(c.server, c.method, c.body)
	return nil
}
