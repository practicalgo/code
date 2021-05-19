// Listing 1.17: chap1/listing6/main.go
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var totalDuration time.Duration = 5

func getName(rd io.Reader, w io.Writer) (string, error) {
	scanner := bufio.NewScanner(rd)
	msg := "Your name please? Press the Enter key when done"
	fmt.Fprintln(w, msg)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("You entered an empty name")
	}
	return name, nil
}

func getNameContext(ctx context.Context, allowedDuration time.Duration) (string, error) {
	var err error
	name := "Default Name"
	c := make(chan error, 1)

	go func() {
		name, err = getName(os.Stdin, os.Stdout)
		c <- err
	}()

	select {
	case <-ctx.Done():
		return name, ctx.Err()
	case err := <-c:
		return name, err
	}
}

func main() {
	allowedDuration := totalDuration * time.Second
	d := time.Now().Add(allowedDuration)

	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	name, err := getNameContext(ctx, allowedDuration)

	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, name)
}
