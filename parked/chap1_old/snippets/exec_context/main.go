package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, "sleep", "20").Run(); err != nil {
		fmt.Fprintln(os.Stdout, err)
	}
}
