package main

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestSetupExecConfig(t *testing.T) {

	type expectedResult struct {
		err error
	}
	type testConfig struct {
		args   []string
		result expectedResult
	}
	tests := []testConfig{
		testConfig{
			args: []string{"-h"},
			result: expectedResult{
				err: fmt.Errorf("flag: help requested"),
			},
		},
		testConfig{
			args: []string{"exec"},
			result: expectedResult{
				err: fmt.Errorf("Invalid command configuration"),
			},
		},
		testConfig{
			args: []string{"-c", "go"},
			result: expectedResult{
				err: nil,
			},
		},
	}
	byteBuf := new(bytes.Buffer)
	w := bufio.NewWriter(byteBuf)

	for _, tc := range tests {
		_, err := setupExecConfig(w, tc.args)
		if tc.result.err != nil && err == nil {
			t.Errorf("Expected non-nil error: %v, got nil error", tc.result.err)
		}
		if tc.result.err == nil && err != nil {
			t.Errorf("Expected nil error, got error: %v", err)
		}
	}
}
