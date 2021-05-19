//Listing 1.3: chap1/listing1/validate_args_test.go
package main

import (
	"errors"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	type expectedResult struct {
		err error
	}
	type testConfig struct {
		c      config
		result expectedResult
	}
	tests := []testConfig{
		testConfig{
			c: config{},
			result: expectedResult{
				err: errors.New("Must specify a number greater than 0"),
			},
		},
		testConfig{
			c: config{numTimes: -1},
			result: expectedResult{
				err: errors.New("Must specify a number greater than 0"),
			},
		},

		testConfig{
			c: config{numTimes: 10},
			result: expectedResult{
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		err := validateArgs(tc.c)
		if tc.result.err != nil && err.Error() != tc.result.err.Error() {
			t.Errorf("Expected error to be: %v, got: %v\n", tc.result.err, err)
		}
		if tc.result.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
	}
}
