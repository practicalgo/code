package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"testing"

	svc "github.com/practicalgo/code/chap9/server-streaming/service"
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
func TestUserService(t *testing.T) {

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
	resp, err := usersClient.GetUser(
		context.Background(),
		&svc.UserGetRequest{Id: "foo-bar", Email: "jane@doe.com"},
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

func TestRepoService(t *testing.T) {

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
	repoClient := svc.NewRepoClient(client)
	stream, err := repoClient.GetRepos(
		context.Background(),
		&svc.RepoGetRequest{CreatorId: "user-123", Id: "repo-123"},
	)
	if err != nil {
		t.Fatal(err)
	}
	var repos []*svc.Repository
	for {
		repo, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		repos = append(repos, repo.Repo)
	}
	if len(repos) != 5 {
		t.Fatalf("Expected to get back 5 repos, got back: %d repos", len(repos))
	}

	for idx, repo := range repos {
		gotRepoName := repo.Name
		expectedRepoName := fmt.Sprintf("repo-%d", idx+1)

		if gotRepoName != expectedRepoName {
			t.Errorf(
				"Expected Repo Name to be: %s, Got: %s",
				expectedRepoName,
				gotRepoName,
			)
		}
	}
}
