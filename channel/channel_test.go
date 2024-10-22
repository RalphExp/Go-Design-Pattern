package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestOrChannel(t *testing.T) {
	// test
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Minute),
		sig(5*time.Minute),
		sig(10*time.Second),
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Hour),
		sig(1*time.Minute),
		sig(2*time.Second),
	)
	duration := time.Since(start)
	if duration > time.Second*5 {
		t.Errorf("execute time too long: %vms", duration.Milliseconds())
	}
}

func TestBridgeChannel(t *testing.T) {
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 100; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	i := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for v := range bridge(ctx, genVals()) {
		fmt.Printf("%v ", v)
		i++
		if i == 50 {
			cancel()
		}
	}
	if i != 50 && i != 51 {
		t.Errorf("unexpected value of loop: %d\n", i)
	}
}
