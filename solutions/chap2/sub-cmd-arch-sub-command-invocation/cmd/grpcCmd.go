// Listing 1.15: chap1/listing5/cmd/grpcCmd.go
package cmd

import (
	"flag"
	"fmt"
	"io"
)

type grpcConfig struct {
	server string
	method string
	body   string
}

func HandleGrpc(w io.Writer, args []string) error {
	c := grpcConfig{}
	fs := flag.NewFlagSet("grpc", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.method, "method", "", "Method to call")
	fs.StringVar(&c.body, "body", "", "Body of request")
	fs.Usage = func() {
		var usageString = `
grpc: A gRPC client.

grpc: <options> server`
		fmt.Fprintf(w, usageString)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}
	c.server = fs.Arg(0)
	fmt.Fprintln(w, "Executing grpc command")
	return nil
}
