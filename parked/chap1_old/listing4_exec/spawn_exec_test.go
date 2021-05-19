package main

import (
	"fmt"
	"regexp"
	"testing"
)

func TestSpawnExec(t *testing.T) {

	type expectedResult struct {
		outputRegexp *regexp.Regexp
		stderrRegexp *regexp.Regexp
		err          error
	}
	type testConfig struct {
		args   []string
		result expectedResult
	}
	tests := []testConfig{
		testConfig{
			args: []string{"version"},
			result: expectedResult{
				err:          nil,
				outputRegexp: regexp.MustCompile("go version go1.[0-9]+"),
			},
		},
		testConfig{
			args: []string{"env", "GOARCH"},
			result: expectedResult{
				err:          nil,
				outputRegexp: regexp.MustCompile("amd64"),
			},
		},
		testConfig{
			args: []string{"arg1"},
			result: expectedResult{
				err:          fmt.Errorf("exit status 2"),
				stderrRegexp: regexp.MustCompile("go arg1: unknown command\nRun 'go help' for usage."),
			},
		},
	}

	c := command{cmd: "go"}

	for _, tc := range tests {
		conf := execConfig{
			cmd:  c,
			args: tc.args,
		}
		result, err := spawnExec(conf)
		if tc.result.err != nil && err == nil {
			t.Errorf("Expected non-nil error: %v, got nil error", tc.result.err)
		}
		if tc.result.err == nil && err != nil {
			t.Errorf("Expected nil error, got error: %v, stderr: %v", err, string(result.stderr))
		}
		if tc.result.outputRegexp != nil {
			matched := tc.result.outputRegexp.MatchString(string(result.stdout))
			if !matched {
				t.Errorf("Stdout: %v did not match the expected regex: %v", result.stdout, tc.result.outputRegexp)
			}
		}
		if tc.result.stderrRegexp != nil {
			matched := tc.result.stderrRegexp.MatchString(string(result.stderr))
			if !matched {
				t.Errorf("Stderr: %v did not match the expected regex: %v", result.stderr, tc.result.stderrRegexp)
			}

		}
	}
}
