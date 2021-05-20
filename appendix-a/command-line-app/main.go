package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"time"

	"os"

	"github.com/practicalgo/code/appendix-a/pkgcli/cmd"
	"github.com/practicalgo/code/appendix-a/pkgcli/config"
	"github.com/practicalgo/code/appendix-a/pkgcli/telemetry"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/rs/zerolog"
)

var errInvalidSubCommand = errors.New("invalid sub-command specified")
var version = "0.1"

type pkgCliInput struct {
	logLevel   int
	statsdAddr string
	jaegerAddr string
}

func printSubCmdUsage(c config.PkgCliConfig, w io.Writer) {
	cmd.HandleRegister(&c, w, []string{"-h"})
	cmd.HandleQuery(&c, w, []string{"-h"})
}

func handleSubCommand(c config.PkgCliConfig, w io.Writer, args []string) error {
	var err error

	if len(args) < 1 {
		err = errInvalidSubCommand
	} else {
		switch args[0] {
		case "register":

			c.Logger = c.Logger.With().Str("command", "register").Logger()
			tStart := time.Now()
			defer func() {
				// We report the metrics in seconds since we are reporting a histogram metric
				// and are using Prometheus as the final metrics server
				duration := time.Since(tStart).Seconds()
				c.Logger.Debug().Msg("Reporting metric")
				c.Metrics.ReportDuration(
					telemetry.DurationMetric{
						Cmd:        "pkgcli.register",
						DurationMs: duration,
						Success:    err == nil,
					},
				)
			}()
			err = cmd.HandleRegister(&c, w, args[1:])
			c.Logger.Debug().Msg("Back from HandleRegister")

		case "query":
			c.Logger = c.Logger.With().Str("command", "query").Logger()
			tStart := time.Now()
			defer func() {
				// We report the metrics in seconds since we are reporting a histogram metric
				// and are using Prometheus as the final metrics server
				duration := time.Since(tStart).Seconds()
				c.Logger.Debug().Msg("Reporting metric")
				c.Metrics.ReportDuration(
					telemetry.DurationMetric{
						Cmd:        "pkgcli.queryr",
						DurationMs: duration,
						Success:    err == nil,
					},
				)
			}()
			err = cmd.HandleQuery(&c, w, args[1:])
		default:
			c.Logger.With().Str("command", args[0])
			err = errInvalidSubCommand
		}
	}
	if errors.Is(err, cmd.ErrNoServerSpecified) || errors.Is(err, errInvalidSubCommand) {
		fmt.Fprintln(w, err)
		printSubCmdUsage(c, w)
	}
	return err
}

func parseArgs(w io.Writer, args []string) (pkgCliInput, []string, error) {
	var logLevel int
	var statsdAddr, jaegerAddr string

	c := pkgCliInput{}
	fs := flag.NewFlagSet("pkgcli", flag.ContinueOnError)
	fs.SetOutput(w)
	origUsage := fs.Usage
	fs.Usage = func() {
		origUsage()
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Environment variables:")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "X_AUTH_TOKEN - Specify auth token")
		printSubCmdUsage(config.PkgCliConfig{}, w)
	}
	fs.IntVar(&logLevel, "log-level", int(zerolog.InfoLevel), "Set log level (-1 to 5), default - 1 INFO logs")
	fs.StringVar(&statsdAddr, "metrics-addr", "", "Metrics server address")
	fs.StringVar(&jaegerAddr, "jaeger-addr", "", "Jaeger server address")

	err := fs.Parse(args)
	if err != nil {
		return c, nil, err
	}
	if !(logLevel >= -1 && logLevel <= 5) {
		fs.Usage()
		return c, nil, errors.New("invalid log level")
	}

	if len(statsdAddr) == 0 {
		fs.Usage()
		return c, nil, errors.New("empty metrics server address")
	}
	if len(jaegerAddr) == 0 {
		fs.Usage()
		return c, nil, errors.New("empty jaeger server address")
	}
	c.logLevel = logLevel
	c.statsdAddr = statsdAddr
	c.jaegerAddr = jaegerAddr
	return c, fs.Args(), nil
}

func main() {
	var tp *sdktrace.TracerProvider
	var cliConfig config.PkgCliConfig

	c, subCmdArgs, err := parseArgs(os.Stdout, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	cliConfig.Logger = telemetry.InitLogging(os.Stdout, version, c.logLevel)

	cliConfig.Metrics, err = telemetry.InitMetrics(c.statsdAddr)
	if err != nil {
		cliConfig.Logger.Fatal().Str("error", err.Error()).Msg(
			"Error initializing metrics system",
		)
	}
	// TODO fix the URL path construction
	cliConfig.Tracer, tp, err = telemetry.InitTracing(
		c.jaegerAddr+"/api/traces", version,
	)
	if err != nil {
		cliConfig.Logger.Fatal().Str("error", err.Error()).Msg(
			"Error initializing tracing system",
		)
	}
	defer func() {
		tp.ForceFlush(context.Background())
		tp.Shutdown(context.Background())
		cliConfig.Metrics.Close()
	}()

	err = handleSubCommand(cliConfig, os.Stdout, subCmdArgs)
	if err != nil {
		// TODO - We have to do this manually here
		// since defer functions are not executed when we do
		// log.Fatal or call os.Exit(1)
		// Likely a good idea to create a new function to avoid code duplication
		// in Line 161
		tp.ForceFlush(context.Background())
		tp.Shutdown(context.Background())
		cliConfig.Metrics.Close()

		cliConfig.Logger.Fatal().Str("error", err.Error()).Msg(
			"Error executing sub-command",
		)
	}
}
