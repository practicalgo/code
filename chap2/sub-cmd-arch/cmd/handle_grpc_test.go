package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func TestHandleGrpc(t *testing.T) {
	usageMessage := `
grpc: A gRPC client.

grpc: <options> server

Options: 
  -body string
    	Body of request
  -method string
    	Method to call
`
	testConfigs := []struct {
		args   []string
		err    error
		output string
	}{
		{
			args: []string{},
			err:  ErrNoServerSpecified,
		},

		{
			args:   []string{"-h"},
			err:    errors.New("flag: help requested"),
			output: usageMessage,
		},

		{
			args:   []string{"-method", "service.host.local/method", "-body", "{}", "http://localhost"},
			err:    nil,
			output: "Executing grpc command\n",
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := HandleGrpc(byteBuf, tc.args)
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
