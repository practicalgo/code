package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type httpConfig struct {
	Url      string
	Verb     string
	JsonData string
}

func handleHttp(w io.Writer, args []string) error {
	var v string
	c := httpConfig{Verb: v}
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.Var(&c, "Verb", "HTTP method")
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.PrintDefaults()
		return errors.New("Url not specified")
	}
	c.Url = fs.Arg(0)
	fmt.Fprintf(w, "Executing http command")
	return nil
}

func printUsage() {
	fmt.Printf("Usage: %s [http|grpc] -h\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "http":
		err := handleHttp(os.Stdout, os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}
	case "grpc":
		fmt.Fprintf(os.Stdout, "Handling grpc command")
	default:
		printUsage()
	}
}
