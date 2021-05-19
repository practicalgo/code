// Listing 1.17: chap1/listing5/handle_command_test.go
package main

import (
	"bytes"
	"testing"
)

func TestHandleCommand(t *testing.T) {
	usageMessage := `Usage: mync [http|grpc] -h

http: A HTTP client.

http: <options> server

Options: 
  -verb string
    	HTTP method (default "GET")

grpc: A gRPC client.

grpc: <options> server

Options: 
  -body string
    	Body of request
  -method string
    	Method to call
`

	type expectedResult struct {
		output string
		err    error
	}
	type testConfig struct {
		args   []string
		result expectedResult
	}

	testConfigs := []testConfig{
		testConfig{
			args: []string{},
			result: expectedResult{
				err:    errInvalidSubCommand,
				output: "Invalid sub-command specified\n" + usageMessage,
			},
		},

		testConfig{
			args: []string{"-h"},
			result: expectedResult{
				err:    nil,
				output: usageMessage,
			},
		},

		testConfig{
			args: []string{"foo"},
			result: expectedResult{
				err:    errInvalidSubCommand,
				output: "Invalid sub-command specified\n" + usageMessage,
			},
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := handleCommand(byteBuf, tc.args)
		if tc.result.err == nil && err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}

		if tc.result.err != nil && err.Error() != tc.result.err.Error() {
			t.Errorf("Expected error %v, got %v", tc.result.err, err)
		}

		if len(tc.result.output) != 0 {
			gotOutput := byteBuf.String()
			if tc.result.output != gotOutput {
				t.Errorf("Expected output to be: %#v, Got: %#v", tc.result.output, gotOutput)
			}
		}
		byteBuf.Reset()
	}
}
