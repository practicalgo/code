package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/practicalgo/code/appendix-b/pkgcli/config"
	"github.com/practicalgo/code/appendix-b/pkgcli/telemetry"
)

func TestQueryCmd(t *testing.T) {

	testConfigs := []struct {
		args   []string
		err    error
		output string
	}{
		{
			args: []string{},
			err:  ErrInvalidQueryArguments,
		},

		{
			args:   []string{"-name", "package1", "-owner", "1", "http://localhost"},
			err:    nil,
			output: "Executing query command\n",
		},
	}

	cliConfig := config.PkgCliConfig{
		Logger: telemetry.InitLogging(os.Stderr, "0.1", 0),
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := HandleQuery(&cliConfig, byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}

		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error %v, got %v", tc.err, err)
		}

		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if tc.output != gotOutput {
				t.Fatalf("Expected output to be: %#v, Got: %#v", tc.output, gotOutput)
			}
		}

		byteBuf.Reset()
	}

}
