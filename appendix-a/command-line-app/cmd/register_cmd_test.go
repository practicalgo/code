package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/practicalgo/code/appendix-a/pkgcli/config"
	"github.com/practicalgo/code/appendix-a/pkgcli/pkgregister"
	"github.com/practicalgo/code/appendix-a/pkgcli/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTestTracer() (telemetry.TraceReporter, *sdktrace.TracerProvider, error) {

	var tr telemetry.TraceReporter

	traceExporter, err := stdouttrace.New()
	if err != nil {
		return tr, nil, err
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

	tr.Client = otel.Tracer("")

	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/baggage/api.md
	version := "0.1-test"
	v1, err := baggage.NewMember("cli_version", version)
	if err != nil {
		return tr, tp, err
	}
	bag, err := baggage.New(v1)
	if err != nil {
		return tr, tp, err
	}
	ctx := context.Background()
	tr.Ctx = baggage.ContextWithBaggage(ctx, bag)

	return tr, tp, nil
}

func testPackageRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Package registration response
		d := pkgregister.PkgRegisterResult{}
		err := r.ParseMultipartForm(5000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mForm := r.MultipartForm
		f := mForm.File["filedata"][0]
		// Construct an artificial package ID to return
		d.ID = fmt.Sprintf("%s-%s-%s", mForm.Value["name"][0], mForm.Value["version"][0], f.Filename)
		jsonData, err := json.Marshal(d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonData))
	} else {
		http.Error(w, "Invalid HTTP method specified", http.StatusMethodNotAllowed)
		return
	}
}

func startTestPackageServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(testPackageRegHandler))
	return ts
}

func TestRegisterCmd(t *testing.T) {

	tmpDir := t.TempDir()
	f, err := os.CreateTemp(tmpDir, "test-package")
	if err != nil {
		t.Fatal(err)
	}

	ts := startTestPackageServer()
	defer ts.Close()

	testConfigs := []struct {
		args   []string
		output string
		err    error
	}{
		{
			args: []string{},
			err:  ErrInvalidRegisterArguments,
		},
		{
			args: []string{"-name", "package1", "-version", "0.1", "-path", "/tmp/package1.tar.gz"},
			err:  ErrNoServerSpecified,
		},

		{
			args: []string{
				"-name", "package1", "-version", "0.1",
				"-path", f.Name(),
				ts.URL,
			},
			err:    nil,
			output: "Uploading package...\n",
		},
	}
	tr, _, err := initTestTracer()
	if err != nil {
		t.Fatal(err)
	}

	mr, _ := telemetry.InitMetrics("127.0.0.1:9125")

	cliConfig := config.PkgCliConfig{
		Logger:  telemetry.InitLogging(os.Stderr, "0.1", 0),
		Tracer:  tr,
		Metrics: mr,
	}
	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := HandleRegister(&cliConfig, byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}

		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error %v, got %v", tc.err, err)
		}

		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if strings.Contains(gotOutput, tc.output) {
				t.Errorf("Expected output to contain: %#v, Got: %#v", tc.output, gotOutput)
			}
		}
		byteBuf.Reset()
	}
}
