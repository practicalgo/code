package handlers

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	users "github.com/practicalgo/code/appendix-b/grpc-server/service"
	"github.com/practicalgo/code/appendix-b/http-server/config"
	"github.com/practicalgo/code/appendix-b/http-server/storage"
	"github.com/practicalgo/code/appendix-b/http-server/telemetry"
	"github.com/practicalgo/code/appendix-b/http-server/testutils"
	"github.com/practicalgo/code/appendix-b/http-server/types"
	_ "gocloud.dev/blob/fileblob"
	"google.golang.org/grpc"
)

func TestPackageGetHandlerNoData(t *testing.T) {
	packageBucket, err := testutils.GetTestBucket(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer packageBucket.Close()

	testC, addr, err := testutils.GetTestDb()
	if err != nil {
		log.Fatal(err)
	}
	defer testC.Terminate(context.Background())
	testDb, err := storage.GetDatabaseConn(
		addr, "package_server",
		"packages_rw", "password",
	)
	if err != nil {
		t.Fatal(err)
	}

	metrics, err := telemetry.InitMetrics("127.0.0.1:9125")
	if err != nil {
		t.Fatal(err)
	}

	tracer, err := testutils.InitTestTracer()
	if err != nil {
		t.Fatal(err)
	}
	s, l := testutils.StartTestGrpcServer()
	defer s.GracefulStop()

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

	c := users.NewUsersClient(conn)

	config := config.AppConfig{
		PackageBucket: packageBucket,
		Db:            testDb,
		UsersSvc:      c,
		Logger:        telemetry.InitLogging(os.Stdout, "0.1", 0),
		Metrics:       metrics,
		Trace:         tracer,
		TraceCtx:      context.Background(),
	}
	ctx, span := config.Trace.Start(config.TraceCtx, "get_handler_test")
	config.Span = span
	config.SpanCtx = ctx
	defer config.Span.End()

	mux := http.NewServeMux()
	SetupHandlers(mux, &config)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := http.Client{}
	resp, err := client.Get(ts.URL + "/api/packages?owner_id=1&name=test-package&version=0.1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
		t.Fatalf(
			"Expected a HTTP 404 response, Got: %v\n",
			resp.StatusCode,
		)
	}
}

func TestPackageGetHandler(t *testing.T) {
	packageBucket, err := testutils.GetTestBucket(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer packageBucket.Close()

	testObjectId := "pkg-0.1-pkg-0.1.tar.gz"

	// create a test object
	err = packageBucket.WriteAll(
		context.Background(),
		testObjectId, []byte("test-data"),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	testC, addr, err := testutils.GetTestDb()
	if err != nil {
		log.Fatal(err)
	}
	defer testC.Terminate(context.Background())
	testDb, err := storage.GetDatabaseConn(
		addr, "package_server",
		"packages_rw", "password",
	)
	if err != nil {
		t.Fatal(err)
	}

	metrics, err := telemetry.InitMetrics("127.0.0.1:9125")
	if err != nil {
		t.Fatal(err)
	}

	tracer, err := testutils.InitTestTracer()
	if err != nil {
		t.Fatal(err)
	}

	s, l := testutils.StartTestGrpcServer()
	defer s.GracefulStop()

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

	c := users.NewUsersClient(conn)

	config := config.AppConfig{
		PackageBucket: packageBucket,
		Db:            testDb,
		UsersSvc:      c,
		Logger:        telemetry.InitLogging(os.Stdout, "0.1", 0),
		Metrics:       metrics,
		Trace:         tracer,
		TraceCtx:      context.Background(),
	}
	ctx, span := config.Trace.Start(config.TraceCtx, "get_handler_test")
	config.Span = span
	config.SpanCtx = ctx
	defer config.Span.End()

	// update package metadata for the test object
	err = storage.UpdateDb(
		context.Background(),
		&config,
		types.PkgRow{
			OwnerId:       1,
			Name:          "pkg",
			Version:       "0.1",
			ObjectStoreId: testObjectId,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	SetupHandlers(mux, &config)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	var redirectUrl string
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirectUrl = req.URL.String()
			return errors.New("no redirect")
		},
	}

	_, err = client.Get(ts.URL + "/api/packages?owner_id=1&name=pkg&version=0.1")
	if err == nil {
		t.Fatal("Expected error: no redirect, Got nil")
	}
	if !strings.HasPrefix(redirectUrl, "file:///") {
		t.Fatalf("Expected redirect url to start with file:///, got: %v", redirectUrl)
	}
}
