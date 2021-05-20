package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/practicalgo/code/appendix-b/http-server/config"
	"github.com/practicalgo/code/appendix-b/http-server/storage"
	"github.com/practicalgo/code/appendix-b/http-server/telemetry"
	"github.com/practicalgo/code/appendix-b/http-server/testutils"
	"github.com/practicalgo/code/appendix-b/http-server/types"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	users "github.com/practicalgo/code/appendix-b/grpc-server/service"
	_ "gocloud.dev/blob/fileblob"

	"google.golang.org/grpc"
)

type pkgData struct {
	Name     string
	Version  string
	Filename string
	Bytes    io.Reader
}

func createMultiPartMessage(data pkgData) ([]byte, string, error) {
	var b bytes.Buffer
	var err error
	var fw io.Writer

	mw := multipart.NewWriter(&b)

	fw, err = mw.CreateFormField("name")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Name)

	fw, err = mw.CreateFormField("version")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Version)

	fw, err = mw.CreateFormFile("filedata", data.Filename)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(fw, data.Bytes)
	if err != nil {
		return nil, "", err
	}
	err = mw.Close()
	if err != nil {
		return nil, "", err
	}

	contentType := mw.FormDataContentType()
	return b.Bytes(), contentType, nil
}

func TestPackageRegHandler(t *testing.T) {
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

	p := pkgData{
		Name:     "mypackage",
		Version:  "0.1",
		Filename: "mypackage-0.1.tar.gz",
		Bytes:    strings.NewReader("data"),
	}

	payload, contentType, err := createMultiPartMessage(p)
	if err != nil {
		t.Fatal(err)
	}
	reader := bytes.NewReader(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ts.URL+"/api/packages", reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-Auth-Token", "test-token")

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	r, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	respData, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	resp := types.PkgRegisterResponse{}
	err = json.Unmarshal(respData, &resp)
	if err != nil {
		t.Fatal(err)
	}
	expectedPackageId := "1/mypackage-0.1-mypackage-0.1.tar.gz"
	if resp.ID != expectedPackageId {
		t.Fatalf(
			"Expected version to be %s, Got: %s\n",
			expectedPackageId,
			resp.ID,
		)
	}
}
