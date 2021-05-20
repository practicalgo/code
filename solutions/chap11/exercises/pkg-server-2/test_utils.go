package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
)

func getTestBucket(tmpDir string) (*blob.Bucket, error) {
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

func getTestDb() (
	testcontainers.Container,
	*sql.DB,
	error,
) {
	bootStrapSqlDir, err := os.Stat("mysql-init")
	if err != nil {
		return nil, nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		if err != nil {
			return nil, nil, err
		}

	}
	bindMountPath := filepath.Join(cwd, bootStrapSqlDir.Name())

	waitForSql := wait.ForSQL("3306/tcp", "mysql",
		func(p nat.Port) string {
			return "root:rootpw@tcp(" +
				"127.0.0.1:" + p.Port() +
				")/package_server"
		})
	waitForSql.WithPollInterval(15 * time.Second)

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
		return mysqlC, nil, err
	}
	addr, err := mysqlC.PortEndpoint(ctx, "3306", "")
	if err != nil {
		return mysqlC, nil, err
	}
	db, err := getDatabaseConn(
		addr, "package_server",
		"packages_rw", "password",
	)
	if err != nil {
		return mysqlC, nil, nil
	}
	return mysqlC, db, nil

}
