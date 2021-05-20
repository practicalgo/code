package testutils

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"

	users "github.com/practicalgo/code/appendix-a/grpc-server/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func GetTestBucket(tmpDir string) (*blob.Bucket, error) {
	myDir, err := os.MkdirTemp(tmpDir, "test-bucket")
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(fmt.Sprintf("file:///%s", myDir))
	if err != nil {
		return nil, err
	}
	opts := fileblob.Options{
		URLSigner: fileblob.NewURLSignerHMAC(
			u,
			[]byte("super secret"),
		),
	}
	return fileblob.OpenBucket(myDir, &opts)
}

func GetTestDb() (testcontainers.Container, string, error) {
	bootStrapSqlDir, err := os.Stat("../mysql-init")
	if err != nil {
		return nil, "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		if err != nil {
			return nil, "", err
		}

	}
	bindMountPath := filepath.Join(cwd+"/../", bootStrapSqlDir.Name())

	waitForSql := wait.ForSQL("3306/tcp", "mysql",
		func(p nat.Port) string {
			return "root:rootpw@tcp(" +
				"127.0.0.1:" + p.Port() +
				")/package_server"
		})
	waitForSql.WithPollInterval(5 * time.Second)

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0.26",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_DATABASE":      "package_server",
			"MYSQL_USER":          "packages_rw",
			"MYSQL_PASSWORD":      "password",
			"MYSQL_ROOT_PASSWORD": "rootpw",
		},
		ImagePlatform: "linux/x86_64",
		BindMounts: map[string]string{
			bindMountPath: "/docker-entrypoint-initdb.d",
		},
		Cmd: []string{
			"--default-authentication-plugin=mysql_native_password",
		},
		WaitingFor: waitForSql,
	}
	ctx := context.Background()
	mysqlC, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		return mysqlC, "", err
	}
	addr, err := mysqlC.PortEndpoint(ctx, "3306", "")
	if err != nil {
		return mysqlC, "", err
	}
	return mysqlC, addr, nil
}

func InitTestTracer() (trace.Tracer, error) {

	traceExporter, err := stdouttrace.New()
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewSimpleSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("PkgServer-Cli"),
			),
		),
	)
	otel.SetTracerProvider(tp)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
	)
	otel.SetTextMapPropagator(propagator)

	return otel.Tracer(""), nil
}

type dummyUserService struct {
	users.UnimplementedUsersServer
}

func (s *dummyUserService) GetUser(
	ctx context.Context,
	in *users.UserGetRequest,
) (*users.UserGetReply, error) {
	u := users.User{
		Id: 1,
	}
	return &users.UserGetReply{User: &u}, nil
}

func StartTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	users.RegisterUsersServer(s, &dummyUserService{})
	go func() {
		err := s.Serve(l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}
