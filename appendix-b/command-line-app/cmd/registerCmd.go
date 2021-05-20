package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/practicalgo/code/appendix-b/pkgcli/config"
	"github.com/practicalgo/code/appendix-b/pkgcli/pkgregister"
)

type PkgRegisterConfig struct {
	name      string
	version   string
	filePath  string
	serverUrl string
}

func validateRegisterConfig(c PkgRegisterConfig) bool {
	if len(c.name) == 0 || len(c.version) == 0 || len(c.filePath) == 0 {
		return false
	}
	return true
}

func HandleRegister(cliConfig *config.PkgCliConfig, w io.Writer, args []string) error {
	c := PkgRegisterConfig{}
	fs := flag.NewFlagSet("register", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.name, "name", "", "Package Name")
	fs.StringVar(&c.version, "version", "", "Package Version")
	fs.StringVar(&c.filePath, "path", "", "Package File Path")

	fs.Usage = func() {
		var usageString = `
register: Upload a package.

regiser: <options> <server URL>`
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

	if !validateRegisterConfig(c) {
		cliConfig.Logger.Debug().Msg("Invalid config")
		return ErrInvalidRegisterArguments
	}

	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}
	c.serverUrl = fs.Arg(0)

	ctx, span := cliConfig.Tracer.Client.Start(
		cliConfig.Tracer.Ctx,
		"pkgquery.register",
	)
	defer span.End()

	cliConfig.Logger.Info().Msg("Uploading package...")
	cliConfig.Logger.Debug().Str("package_name", c.name).Str("package_version", c.version).Str("server_url", c.serverUrl)
	f, err := os.Open(c.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	pkgData := pkgregister.PkgData{
		Name:     c.name,
		Version:  c.version,
		Filename: c.filePath,
		Bytes:    f,
	}

	cliConfig.Logger.Debug().Msg("Making HTTP POST request")
	cliConfig.Logger.Debug().Msg(fmt.Sprintf("%#v", pkgData))

	resp, err := pkgregister.RegisterPackage(ctx, cliConfig, c.serverUrl, pkgData)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Package uploaded: %s", resp.ID)
	return nil
}
