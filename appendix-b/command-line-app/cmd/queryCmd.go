package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/practicalgo/code/appendix-b/pkgcli/config"
)

type PkgQueryConfig struct {
	name      string
	version   string
	owner     string
	serverUrl string
}

func validateQueryConfig(c PkgQueryConfig) bool {
	if len(c.name) == 0 && len(c.owner) == 0 {
		return false
	}
	if len(c.version) != 0 && len(c.name) == 0 {
		return false
	}
	return true
}

func HandleQuery(pkgCliConfig *config.PkgCliConfig, w io.Writer, args []string) error {
	c := PkgQueryConfig{}
	fs := flag.NewFlagSet("query", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.name, "name", "", "Name of package")
	fs.StringVar(&c.version, "version", "", "Version of package")
	fs.StringVar(&c.owner, "owner", "", "Owner of package")
	fs.Usage = func() {
		var usageString = `
query: Query Packages.

query: <options> serverUrl`
		fmt.Fprint(w, usageString)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}
	if !validateQueryConfig(c) {
		return ErrInvalidQueryArguments
	}

	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}
	c.serverUrl = fs.Arg(0)
	// TODO implement this
	fmt.Fprintln(w, "Executing query command")
	return nil
}
