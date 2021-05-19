package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type execConfig struct {
	cmd  command
	args []string
}

func (c execConfig) valid() bool {
	if len(c.cmd.cmd) == 0 {
		return false
	}

	if len(c.args) > 5 {
		return false
	}
	return true
}

type execOutput struct {
	stdout   []byte
	stderr   []byte
	exitcode int
	err      error
}

type command struct {
	cmd string
}

func (c *command) Set(value string) error {
	allowedCommands := []string{"go", "sleep"}
	allowed := false
	for _, c := range allowedCommands {
		if value == c {
			allowed = true
		}
	}
	if !allowed {
		return errors.New(fmt.Sprintf("Command not allowed: %s", value))
	}
	c.cmd = value
	return nil
}

func (c *command) String() string {
	return c.cmd
}

func setupExecConfig(w io.Writer, args []string) (execConfig, error) {
	c := execConfig{}
	fs := flag.NewFlagSet("exec", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.Var(&c.cmd, "c", "A command to execute")
	err := fs.Parse(args)
	if err != nil {
		return c, err
	}
	c.args = fs.Args()
	if !c.valid() {
		return c, fmt.Errorf("Invalid command configuration")
	}
	return c, nil
}

func spawnExec(ctx context.Context, c execConfig) (execOutput, error) {
	var result execOutput
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, c.cmd.cmd, c.args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	result.err = err
	if e, ok := err.(*exec.ExitError); ok {
		result.exitcode = e.ExitCode()
	} else {
		result.exitcode = -1
	}
	result.stdout = stdout.Bytes()
	result.stderr = stderr.Bytes()
	return result, err
}

func handleExec(w io.Writer, args []string) (execOutput, error) {
	c, err := setupExecConfig(w, args)
	if err != nil {
		return execOutput{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	o, err := spawnExec(ctx, c)
	return o, err
}

func printUsage() {
	fmt.Printf("Usage: %s exec\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "exec":
		result, err := handleExec(os.Stdout, os.Args[2:])
		if err != nil {
			fmt.Printf("Error executing exec: %v\n", err)
			if len(result.stderr) != 0 {
				fmt.Printf("Stderr: %s\n", result.stderr)
			}
		} else {
			fmt.Printf("%s", result.stdout)
		}
	default:
		printUsage()
	}
}
