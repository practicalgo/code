//Listing 1.2: chap1/listing1/parse_args_test.go
package main

import (
	"errors"
	"testing"
)

func TestParseArgs(t *testing.T) {
	type expectedResult struct {
		err error
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
				err:    nil,
				config: config{printUsage: true, numTimes: 0},
			},
		},
		testConfig{
			args: []string{"10"},
			result: expectedResult{
				err:    nil,
				config: config{printUsage: false, numTimes: 10},
			},
		},
		testConfig{
			args: []string{"abc"},
			result: expectedResult{
				err:    errors.New("strconv.Atoi: parsing \"abc\": invalid syntax"),
				config: config{printUsage: false, numTimes: 0},
			},
		},
		testConfig{
			args: []string{"1", "foo"},
			result: expectedResult{
				err:    errors.New("Invalid number of arguments"),
				config: config{printUsage: false, numTimes: 0},
			},
		},
	}

	for _, tc := range tests {
		c, err := parseArgs(tc.args)
		if tc.result.err != nil && err.Error() != tc.result.err.Error() {
			t.Errorf("Expected error to be: %v, got: %v\n", tc.result.err, err)
		}
		if tc.result.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
		if c.printUsage != tc.result.printUsage {
			t.Errorf("Expected printUsage to be: %v, got: %v\n", tc.result.printUsage, c.printUsage)
		}
		if c.numTimes != tc.result.numTimes {
			t.Errorf("Expected numTimes to be: %v, got: %v\n", tc.result.numTimes, c.numTimes)
		}
	}
}
