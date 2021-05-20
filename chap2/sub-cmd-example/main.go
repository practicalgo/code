package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func handleCmdA(w io.Writer, args []string) error {
	var v string
	fs := flag.NewFlagSet("cmd-a", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&v, "verb", "argument-value", "Argument 1")
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Executing command A")
	return nil
}

func handleCmdB(w io.Writer, args []string) error {
	var v string
	fs := flag.NewFlagSet("cmd-b", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&v, "verb", "argument-value", "Argument 1")
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Executing command B")
	return nil
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [cmd-a|cmd-b] -h\n", os.Args[0])
	handleCmdA(w, []string{"-h"})
	handleCmdB(w, []string{"-h"})

}

func main() {
	var err error
	if len(os.Args) < 2 {
		printUsage(os.Stdout)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "cmd-a":
		err = handleCmdA(os.Stdout, os.Args[2:])
	case "cmd-b":
		err = handleCmdB(os.Stdout, os.Args[2:])
	default:
		printUsage(os.Stdout)
	}

	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
