//Listing 1.9: chap1/listing3/parse_args_test.go
package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestParseArgs(t *testing.T) {
	type expectedResult struct {
		output string
		err    error
		config
	}
	type testConfig struct {
		args   []string
		result expectedResult
	}
	tests := []testConfig{
		testConfig{
			args: []string{"-h"},
			result: expectedResult{
				output: `
A greeter application which prints the name you entered a specified number of times.

Usage of greeter: <options> [name]

Options: 
  -n int
    	Number of times to greet
`,
				err:    errors.New("flag: help requested"),
				config: config{numTimes: 0},
			},
		},
		testConfig{
			args: []string{"-n", "10"},
			result: expectedResult{
				err:    nil,
				config: config{numTimes: 10},
			},
		},
		testConfig{
			args: []string{"-n", "abc"},
			result: expectedResult{
				err:    errors.New("invalid value \"abc\" for flag -n: parse error"),
				config: config{numTimes: 0},
			},
		},
		testConfig{
			args: []string{"-n", "1", "John Doe"},
			result: expectedResult{
				err:    nil,
				config: config{numTimes: 1, name: "John Doe"},
			},
		},
		testConfig{
			args: []string{"-n", "1", "John", "Doe"},
			result: expectedResult{
				err:    errors.New("More than one positional argument specified"),
				config: config{numTimes: 1},
			},
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		c, err := parseArgs(byteBuf, tc.args)
		if tc.result.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
		if tc.result.err != nil && err.Error() != tc.result.err.Error() {
			t.Errorf("Expected error to be: %v, got: %v\n", tc.result.err, err)
		}
		if c.numTimes != tc.result.numTimes {
			t.Errorf("Expected numTimes to be: %v, got: %v\n", tc.result.numTimes, c.numTimes)
		}
		gotMsg := byteBuf.String()
		if len(tc.result.output) != 0 && gotMsg != tc.result.output {
			t.Errorf("Expected stdout message to be: %#v, Got: %#v\n", tc.result.output, gotMsg)
		}
		byteBuf.Reset()
	}
}
