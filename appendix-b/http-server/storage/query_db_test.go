package storage

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/practicalgo/code/appendix-b/http-server/config"
	"github.com/practicalgo/code/appendix-b/http-server/telemetry"
	"github.com/practicalgo/code/appendix-b/http-server/testutils"
	"github.com/practicalgo/code/appendix-b/http-server/types"
)

func TestQueryDb(t *testing.T) {

	testC, addr, err := testutils.GetTestDb()
	if err != nil {
		log.Fatal(err)
	}
	defer testC.Terminate(context.Background())
	testDb, err := GetDatabaseConn(
		addr, "package_server",
		"packages_rw", "password",
	)
	if err != nil {
		t.Fatal(err)
	}

	tracer, err := testutils.InitTestTracer()
	if err != nil {
		t.Fatal(err)
	}

	config := config.AppConfig{
		Logger:   telemetry.InitLogging(os.Stdout, "0.1", 0),
		Db:       testDb,
		Trace:    tracer,
		TraceCtx: context.Background(),
	}
	ctx, span := config.Trace.Start(config.TraceCtx, "query_db_test")
	config.Span = span
	config.SpanCtx = ctx
	defer config.Span.End()

	// update package metadata for the test object
	err = UpdateDb(
		context.Background(),
		&config,
		types.PkgRow{
			OwnerId:       1,
			Name:          "pkg",
			Version:       "0.2",
			ObjectStoreId: "pkg-0.2-pkg-0.2.tar.gz",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// update package metadata for the test object
	err = UpdateDb(
		context.Background(),
		&config,
		types.PkgRow{
			OwnerId:       2,
			Name:          "pkg",
			Version:       "0.3",
			ObjectStoreId: "pkg-0.3-pkg-0.3.tar.gz",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	results, err := QueryDb(
		&config,
		types.PkgQueryParams{
			OwnerId: 2,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"Expected: 1 row, Got: %d", len(results),
		)
	}

	layout := "2006-01-02 15:04:05"
	created := results[0].Created
	parsedTime, err := time.Parse(layout, created)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", parsedTime.Local().String())
}
