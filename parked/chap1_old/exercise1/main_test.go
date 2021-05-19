package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	cmd := exec.Command("go", "build", "-o", "./test_application")
	if _, err := cmd.Output(); err != nil {
		log.Fatalf("Error building application: %v\n", err)
	}
	os.Exit(m.Run())
}

func TestMainE2E(t *testing.T) {
	type testConfig struct {
		arg      string
		input    string
		output   string
		exitCode int
	}

	testConfigs := []testConfig{
		testConfig{arg: "-h", exitCode: 1},
		testConfig{arg: "a", exitCode: 1},
		testConfig{arg: "1", input: "Jane C", output: "Your name please? Press the Enter key when done.\nNice to meet you Jane C\n"},
	}

	for _, tc := range testConfigs {
		cmd := exec.Command("./test_application", tc.arg)
		if len(tc.input) != 0 {
			cmd.Stdin = strings.NewReader(tc.input)
		}
		var output []byte
		var err error
		if output, err = cmd.Output(); err != nil {
			if e, ok := err.(*exec.ExitError); ok {
				if e.ExitCode() != tc.exitCode {
					t.Errorf("Expected exit code to be %v, Got: %v", tc.exitCode, e.ExitCode())
				}
			}
		} else {
			if tc.output != string(output) {
				t.Errorf("Expected output to be: %#v, Got: %#v\n", tc.output, string(output))
			}
		}
	}
}
