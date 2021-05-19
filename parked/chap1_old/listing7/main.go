// Listing 1.18: chap1/listing7/main.go
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func ctxCancellableTimeout() (context.Context, context.CancelFunc, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ctxS, sCancel := context.WithCancel(ctx)
	return ctxS, cancel, sCancel
}

func setupSignalHandler(w io.Writer, cancelFunc context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		s := <-c
		fmt.Fprintf(w, "Got signal:%v\n", s)
		cancelFunc()
	}()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stdout, "Usage: %s <command> <argument>\n", os.Args[0])
		os.Exit(1)
	}
	command := os.Args[1]
	arg := os.Args[2]

	ctx, cancel, sCancel := ctxCancellableTimeout()
	defer cancel()
	defer sCancel()
	setupSignalHandler(os.Stdout, sCancel)

	if err := exec.CommandContext(ctx, command, arg).Run(); err != nil {
		fmt.Fprintln(os.Stdout, err)
	}
}
