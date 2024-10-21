package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			// case <-create_or_channel(append(channels[3:], orDone)...):
			case <-or(channels[3:]...):
			}
		}
	}()
	return orDone
}

func main() {
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
	fmt.Printf("done after %v", time.Since(start))
}
