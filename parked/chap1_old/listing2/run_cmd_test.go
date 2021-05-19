package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	type testConfig struct {
		c              config
		input          string
		expectedOutput string
		err            error
	}

	tests := []testConfig{
		testConfig{
			c:              config{numTimes: 5},
			input:          "",
			expectedOutput: strings.Repeat("Your name please? Press the Enter key when done.\n", 1),
			err:            errors.New("You didn't enter your name"),
		},
		testConfig{
			c:              config{numTimes: 5},
			input:          "Bill Bryson",
			expectedOutput: "Your name please? Press the Enter key when done.\n" + strings.Repeat("Nice to meet you Bill Bryson\n", 5),
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		rd := strings.NewReader(tc.input)
		err := runCmd(rd, byteBuf, tc.c)
		if err != nil && tc.err == nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
		if tc.err != nil {
			if err.Error() != tc.err.Error() {
				t.Errorf("Expected error: %v, Got error: %v\n", tc.err.Error(), err.Error())
			}
		}
		gotMsg := byteBuf.String()
		if gotMsg != tc.expectedOutput {
			t.Errorf("Expected stdout message to be: %v, Got: %v\n", tc.expectedOutput, gotMsg)
		}
		byteBuf.Reset()
	}
}
