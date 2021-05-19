// Listing 1.13: chap1/listing5/cmd/handle_http_test.go
package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func TestHandleHttp(t *testing.T) {
	usageMessage := `
http: A HTTP client.

http: <options> server

Options: 
  -verb string
    	HTTP method (default "GET")
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
				err: ErrNoServerSpecified,
			},
		},

		testConfig{
			args: []string{"-h"},
			result: expectedResult{
				err:    errors.New("flag: help requested"),
				output: usageMessage,
			},
		},

		testConfig{
			args: []string{"http://localhost"},
			result: expectedResult{
				err:    nil,
				output: "Executing http command\n",
			},
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := HandleHttp(byteBuf, tc.args)
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
