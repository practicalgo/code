package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	allowedDuration := 5 * time.Second
	d := time.Now().Add(allowedDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	fmt.Printf("%#v", ctx)
}
