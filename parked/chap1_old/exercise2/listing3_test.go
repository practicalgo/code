package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetupFlagSet(t *testing.T) {
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
				err: errors.New("flag: help requested"),
			},
		},
		testConfig{
			args: []string{"-o", "index.html"},
			result: expectedResult{
				err: nil,
			},
		},
		testConfig{
			args: []string{"-o", "index.html", "foo"},
			result: expectedResult{
				err: errors.New("Invalid command line argument"),
			},
		},
	}

	byteBuf := new(bytes.Buffer)
	w := bufio.NewWriter(byteBuf)
	for _, tc := range tests {
		_, err := setupFlagSet(w, tc.args)
		if tc.result.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
		if tc.result.err != nil && err.Error() != tc.result.err.Error() {
			t.Errorf("Expected error to be: %v, got: %v\n", tc.result.err, err)
		}
		if tc.result.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
	}
}

func TestValidateArgs(t *testing.T) {
	type expectedResult struct {
		err error
	}
	type testConfig struct {
		args   []string
		result expectedResult
	}
	tests := []testConfig{
		testConfig{
			args: []string{},
			result: expectedResult{
				err: errors.New("Must specify -o"),
			},
		},
		testConfig{
			args: []string{"-o", "index.html"},
			result: expectedResult{
				err: nil,
			},
		},
	}

	byteBuf := new(bytes.Buffer)
	w := bufio.NewWriter(byteBuf)
	for _, tc := range tests {
		c, err := setupFlagSet(w, tc.args)
		err = validateArgs(c)
		if tc.result.err != nil && err.Error() != tc.result.err.Error() {
			t.Errorf("Expected error to be: %v, got: %v\n", tc.result.err, err)
		}
		if tc.result.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
	}
}

func TestRunCmd(t *testing.T) {
	type testConfig struct {
		c                    config
		input                string
		expectedOutput       string
		expectedFileContents string
		err                  error
	}

	tempDir := t.TempDir()
	testOutputPath := filepath.Join(tempDir, "index.html")

	tests := []testConfig{
		testConfig{
			c:              config{outputPath: filepath.Join(tempDir, "index.html")},
			input:          "",
			expectedOutput: strings.Repeat("Your name please? Press the Enter key when done.\n", 1),
			err:            errors.New("Error reading name: You didn't enter your name"),
		},
		testConfig{
			c:                    config{outputPath: testOutputPath},
			input:                "Bill Bryson",
			expectedOutput:       fmt.Sprintf("Your name please? Press the Enter key when done.\nYour webpage has been created in: %s\n", testOutputPath),
			expectedFileContents: "<h1> Hello world, welcome to Bill Bryson's website.</h1>",
		},
	}

	for _, tc := range tests {
		rd := bufio.NewReader(strings.NewReader(tc.input))
		byteBuf := new(bytes.Buffer)
		w := bufio.NewWriter(byteBuf)
		f, err := os.Create(tc.c.outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating webpage: %s\n", err.Error())
			os.Exit(1)
		}
		defer f.Close()

		err = runCmd(rd, w, f, tc.c)
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

		if len(tc.expectedFileContents) != 0 {
			fileContents, err := ioutil.ReadFile(tc.c.outputPath)
			if err != nil {
				t.Errorf("Error reading generated file: %w", err)
			}
			gotFileContents := string(fileContents)
			if gotFileContents != tc.expectedFileContents {
				t.Errorf("Expected file contents to be: %v, Got: %v\n", tc.expectedFileContents, gotFileContents)
			}
		}
	}
}
