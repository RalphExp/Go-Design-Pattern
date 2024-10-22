package main

import (
	"context"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func RandIntStream(ctx context.Context) <-chan interface{} {
	return repeatFn(ctx, func() interface{} {
		return rand.Intn(10000000000)
	})
}

func PrimeFinder(ctx context.Context, src <-chan interface{}) <-chan interface{} {
	return pipe(ctx, src, checkPrime)
}

func PrimeFinder2(ctx context.Context, src <-chan interface{}) <-chan interface{} {
	numFinders := runtime.NumCPU()
	finders := make([]<-chan interface{}, numFinders)
	for i := 0; i < numFinders; i++ {
		finders[i] = PrimeFinder(ctx, src)
	}
	return fanIn(ctx, finders...)
}

func TestPrimeFinder(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()
	t.Log("Primes:")
	for prime := range take(ctx, PrimeFinder(ctx, RandIntStream(ctx)), 30) {
		t.Logf("\t%d\n", prime)
	}
	t.Logf("Search took: %v", time.Since(start))
}

func TestPrimeFinder2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()
	t.Log("Primes:")
	c := PrimeFinder2(ctx, RandIntStream(ctx))
	for prime := range take(ctx, c, 30) {
		t.Logf("\t%d\n", prime)
	}
	t.Logf("Search took: %v", time.Since(start))
}
