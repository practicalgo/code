package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type Cmd struct {
	Name     string
	Desc     string
	Func     func(io.Writer, []string) error
	Commands map[string]*Cmd
}

type httpConfig struct {
	Url  string
	Verb string
}

var errInvalidPosArgSpecified = errors.New("Only one positional argument specified")

func handleHttp(w io.Writer, args []string) error {
	var v string
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&v, "verb", "GET", "HTTP method")
	c := httpConfig{Verb: v}
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return errInvalidPosArgSpecified
	}
	c.Url = fs.Arg(0)
	fmt.Fprintf(w, "Executing http command")
	return nil
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [http|grpc] -h\n", os.Args[0])
}

func main() {
	var err error
	if len(os.Args) < 2 {
		printUsage(os.Stdout)
		os.Exit(1)
	}
	mainCmd := Cmd{
		Name: "myclient",
	}
	httpCmd := Cmd{
		Name: "http",
		Func: handleHttp,
	}
	mainCmd.Commands[httpCmd.Name] = httpCmd

	if _, ok := mainCmd.Command[os.Args[1]]; !ok {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "http":
		err = handleHttp(os.Stdout, os.Args[2:])
	case "grpc":
		fmt.Fprintf(os.Stdout, "Handling grpc command")
	default:
		printUsage(os.Stdout)
	}

	if err != nil {
		if errors.Is(err, errInvalidPosArgSpecified) {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}
