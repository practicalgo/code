package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"os"

	"github.com/practicalgo/code/appendix-b/pkgcli/cmd"
	"github.com/practicalgo/code/appendix-b/pkgcli/config"
	"github.com/practicalgo/code/appendix-b/pkgcli/telemetry"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	uberconfig "go.uber.org/config"
)

var errInvalidSubCommand = errors.New("invalid sub-command specified")
var version = "0.1"

type serverCfg struct {
	AuthToken string `yaml:"auth_token"`
}

type telemetryCfg struct {
	LogLevel   int    `yaml:"log_level"`
	StatsdAddr string `yaml:"statsd_addr"`
	JaegerAddr string `yaml:"jaeger_addr"`
}

type pkgCliInput struct {
	Server    serverCfg
	Telemetry telemetryCfg
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

func parseArgs(w io.Writer, args []string, configFilePath string) ([]string, error) {
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

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	return fs.Args(), nil
}

func readConfig(configFilePath string) (*pkgCliInput, error) {
	c := pkgCliInput{}

	provider, err := uberconfig.NewYAML(
		uberconfig.File(configFilePath),
		uberconfig.Expand(os.LookupEnv),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := provider.Get(uberconfig.Root).Populate(&c); err != nil {
		return nil, err
	}

	if len(c.Server.AuthToken) == 0 {
		return nil, err
	}

	if c.Telemetry.LogLevel < -1 || c.Telemetry.LogLevel > 5 {
		return nil, errors.New("invalid log level")
	}

	if len(c.Telemetry.StatsdAddr) == 0 {
		return nil, errors.New("empty statsd address specified")
	}
	if len(c.Telemetry.JaegerAddr) == 0 {
		return nil, errors.New("empty jaeger address specified")
	}

	return &c, nil
}

func main() {
	var tp *sdktrace.TracerProvider
	var cliConfig config.PkgCliConfig

	configFilePath := os.Getenv("CONFIG_FILE")
	if len(configFilePath) == 0 {
		configFilePath = "config.yml"
	}
	subCmdArgs, err := parseArgs(os.Stdout, os.Args[1:], configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	c, err := readConfig(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	cliConfig.Token = c.Server.AuthToken
	cliConfig.Logger = telemetry.InitLogging(os.Stdout, version, c.Telemetry.LogLevel)
	cliConfig.Metrics, err = telemetry.InitMetrics(c.Telemetry.StatsdAddr)
	if err != nil {
		cliConfig.Logger.Fatal().Str("error", err.Error()).Msg(
			"Error initializing metrics system",
		)
	}
	cliConfig.Tracer, tp, err = telemetry.InitTracing(
		c.Telemetry.JaegerAddr+"/api/traces", version,
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

	cliConfig.Logger.Debug().Msgf("%#v", cliConfig)

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
