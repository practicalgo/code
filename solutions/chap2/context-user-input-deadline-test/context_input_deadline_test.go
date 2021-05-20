package main

import (
	"context"
	"testing"
	"time"
)

func TestTroubleFunc(t *testing.T) {
	allowedDuration := 500 * time.Second
	var cancelled []error
	numIterations := 10
	for i := 1; i <= numIterations; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), allowedDuration)
		defer cancel()
		_, err := getNameContext(ctx, allowedDuration)
		cancelled = append(cancelled, err)
	}
	contextExceeded := 0
	for _, e := range cancelled {
		if e != nil {
			contextExceeded++
		}
	}

	if contextExceeded != numIterations/2 {
		t.Errorf("Expected context deadline to exceed %v times, Got: %v", numIterations/2, contextExceeded)
	}
}
