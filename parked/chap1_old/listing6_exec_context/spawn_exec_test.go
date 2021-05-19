package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestSpawnExec(t *testing.T) {

	type expectedResult struct {
		outputRegexp *regexp.Regexp
		stderrRegexp *regexp.Regexp
		err          error
	}
	type testConfig struct {
		c      command
		ctx    context.Context
		args   []string
		result expectedResult
	}
	ctxTimeoutHundredMs, cancel1 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel1()

	ctxTimeoutTens, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	tests := []testConfig{
		testConfig{
			c:    command{cmd: "go"},
			args: []string{"version"},
			result: expectedResult{
				err:          nil,
				outputRegexp: regexp.MustCompile("go version go1.[0-9]+"),
			},
		},
		testConfig{
			c:    command{cmd: "go"},
			args: []string{"env", "GOARCH"},
			result: expectedResult{
				err:          nil,
				outputRegexp: regexp.MustCompile("amd64"),
			},
		},
		testConfig{
			c:    command{cmd: "go"},
			args: []string{"arg1"},
			result: expectedResult{
				err:          fmt.Errorf("exit status 2"),
				stderrRegexp: regexp.MustCompile("go arg1: unknown command\nRun 'go help' for usage."),
			},
		},
		testConfig{
			c:    command{cmd: "sleep"},
			ctx:  ctxTimeoutHundredMs,
			args: []string{"1"},
			result: expectedResult{
				err: fmt.Errorf("signal: killed"),
			},
		},
		testConfig{
			c:    command{cmd: "sleep"},
			ctx:  ctxTimeoutTens,
			args: []string{"1"},
			result: expectedResult{
				err: nil,
			},
		},
	}

	var ctx context.Context
	for _, tc := range tests {
		conf := execConfig{
			cmd:  tc.c,
			args: tc.args,
		}
		if tc.ctx == nil {
			ctx = context.Background()
		} else {
			ctx = tc.ctx
		}
		result, err := spawnExec(ctx, conf)
		if tc.result.err != nil {
			if strings.Index(err.Error(), tc.result.err.Error()) == -1 {
				t.Errorf("Expected non-nil error to be: %v, got error: %v", tc.result.err, err)
			}
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
				t.Errorf("Stderr: %v (Error: %v) did not match the expected regex: %v", result.stderr, err, tc.result.stderrRegexp)
			}

		}
	}
}
