package main

import (
	"context"
	"testing"
	"time"
)

func TestPrimeFinder(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()
	t.Log("Primes:")
	for prime := range take(ctx, PrimeFinder(ctx, RandIntStream(ctx)), 10) {
		t.Logf("\t%d\n", prime)
	}
	t.Logf("Search took: %v", time.Since(start))
}
