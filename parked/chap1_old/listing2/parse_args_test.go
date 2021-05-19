//Listing 1.6: chap1/listing3/parse_args_test.go
package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestParseArgs(t *testing.T) {
	type expectedResult struct {
		err      error
		numTimes int
	}
	type testConfig struct {
		args   []string
		result expectedResult
	}
	tests := []testConfig{
		testConfig{
			args: []string{"-h"},
			result: expectedResult{
				err:      errors.New("flag: help requested"),
				numTimes: 0,
			},
		},
		testConfig{
			args: []string{"-n", "10"},
			result: expectedResult{
				err:      nil,
				numTimes: 10,
			},
		},
		testConfig{
			args: []string{"-n", "abc"},
			result: expectedResult{
				err:      errors.New("invalid value \"abc\" for flag -n: parse error"),
				numTimes: 0,
			},
		},
		testConfig{
			args: []string{"-n", "1", "foo"},
			result: expectedResult{
				err:      errors.New("Positional arguments specified"),
				numTimes: 1,
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
		byteBuf.Reset()
	}
}
