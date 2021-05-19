package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type config struct {
	outputPath string
}

var usageString = fmt.Sprintf(`Usage: %s <integer> [-h|--help]

A greeter application which prints the name you entered <integer> number of times.
`, os.Args[0])

func writeString(w io.Writer, msg string) error {
	writer := bufio.NewWriter(w)
	n, err := writer.WriteString(msg)
	if n != len(msg) {
		return err
	}
	err = writer.Flush()
	return err
}

func getName(rd io.Reader, w io.Writer) (string, error) {
	scanner := bufio.NewScanner(rd)
	msg := "Your name please? Press the Enter key when done.\n"
	err := writeString(w, msg)
	if err != nil {
		return "", fmt.Errorf("Error displaying input prompt: %w", err)
	}

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("You didn't enter your name")
	}
	return name, nil
}

func generateWebpage(c config, name string, w io.Writer) error {
	htmlTemplate := `<h1> Hello world, welcome to %s's website.</h1>`
	msg := fmt.Sprintf(htmlTemplate, name)
	err := writeString(w, msg)
	if err != nil {
		return err
	}
	return nil
}

func runCmd(rd io.Reader, inputW io.Writer, outputW io.Writer, c config) error {
	name, err := getName(rd, inputW)
	if err != nil {
		return fmt.Errorf("Error reading name: %w", err)
	}

	err = generateWebpage(c, name, outputW)
	if err != nil {
		return fmt.Errorf("Error creating webpage: %w", err)
	}

	msg := fmt.Sprintf("Your webpage has been created in: %s\n", c.outputPath)
	err = writeString(inputW, msg)
	if err != nil {
		return fmt.Errorf("Error showing success message: %w", err)
	}
	return nil
}

func validateArgs(c config) error {
	if len(c.outputPath) == 0 {
		return errors.New("Must specify -o")
	}
	return nil
}

func setupFlagSet(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("webpage-generator", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.outputPath, "o", "", "Generate your web page at this path")
	err := fs.Parse(args)
	if err != nil {
		return c, err
	}
	if fs.NArg() != 0 {
		return c, errors.New("Invalid command line argument")
	}
	return c, nil
}

func main() {
	c, err := setupFlagSet(os.Stderr, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
	err = validateArgs(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	f, err := os.Create(c.outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating webpage: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()
	err = runCmd(os.Stdin, os.Stdout, f, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
}
