package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	svc "github.com/practicalgo/code/chap10/client-resiliency/service"
	users "github.com/practicalgo/code/chap10/client-resiliency/service"
	"google.golang.org/grpc"
)

func setupGrpcConn(addr string) (*grpc.ClientConn, context.CancelFunc, error) {
	log.Printf("Connecting to server on %s\n", addr)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithReturnConnectionError(),
	)
	return conn, cancel, err
}

func getUserServiceClient(conn *grpc.ClientConn) svc.UsersClient {
	return svc.NewUsersClient(conn)
}

func getUser(
	client svc.UsersClient,
	u *svc.UserGetRequest,
) (*svc.UserGetReply, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	resp, err := client.GetUser(
		ctx,
		u,
		grpc.WaitForReady(true),
	)
	return resp, cancel, err
}

func createHelpStream(c svc.UsersClient) (
	users.Users_GetHelpClient, error,
) {
	return c.GetHelp(
		context.Background(),
		grpc.WaitForReady(true),
	)
}

func setupChat(r io.Reader, w io.Writer, c svc.UsersClient) (err error) {

	var clientConn = make(chan svc.Users_GetHelpClient)
	var done = make(chan bool)

	stream, err := createHelpStream(c)
	defer stream.CloseSend()
	if err != nil {
		return err
	}

	go func() {
		for {
			clientConn <- stream
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
			}
			if err != nil {
				log.Printf("Recreating stream.")
				stream, err = createHelpStream(c)
				if err != nil {
					close(clientConn)
					done <- true
				}
			} else {
				fmt.Printf("Response: %s\n", resp.Response)
				if resp.Response == "hello-10" {
					done <- true
				}
			}
		}
	}()

	requestMsg := "hello"
	msgCount := 1
	for {
		if msgCount > 10 {
			break
		}
		stream = <-clientConn
		if stream == nil {
			break
		}
		request := svc.UserHelpRequest{
			Request: fmt.Sprintf("%s-%d", requestMsg, msgCount),
		}
		err := stream.Send(&request)
		if err != nil {
			log.Printf("Send error: %v. Will retry.\n", err)
		} else {
			log.Printf("Request sent: %d\n", msgCount)
			msgCount += 1
		}
	}

	<-done
	return stream.CloseSend()
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal(
			"Specify a gRPC server and method to call",
		)
	}
	serverAddr := os.Args[1]
	methodName := os.Args[2]

	if methodName == "GetUser" && len(os.Args) != 4 {
		log.Fatal("Specify an email address for the user")
	}

	conn, cancel, err := setupGrpcConn(serverAddr)
	if err != nil {
		log.Fatalf("Error in creating connection: %v\n", err)
	}

	defer conn.Close()
	defer cancel()

	c := getUserServiceClient(conn)

	switch methodName {
	case "GetUser":
		for i := 1; i <= 5; i++ {
			log.Printf("Request: %d\n", i)
			userEmail := os.Args[3]
			result, cancel, err := getUser(
				c,
				&svc.UserGetRequest{Email: userEmail},
			)
			defer cancel()
			if err != nil {
				log.Fatalf("getUser failed: %v", err)
			}
			fmt.Fprintf(
				os.Stdout,
				"User: %s %s\n",
				result.User.FirstName,
				result.User.LastName,
			)
			log.Printf("Going to sleep for 1 minute")
			time.Sleep(60 * time.Second)
		}
	case "GetHelp":
		err = setupChat(os.Stdin, os.Stdout, c)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Unrecognized method name")
	}
}
