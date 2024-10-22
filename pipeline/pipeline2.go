package main

import (
	"context"
	"sync"
)

func repeat(
	ctx context.Context,
	values ...interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-ctx.Done():
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

func take(
	ctx context.Context,
	valueStream <-chan interface{},
	num int,
) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-ctx.Done():
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

func repeatFn(
	ctx context.Context,
	fn func() interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-ctx.Done():
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

func pipe(
	ctx context.Context,
	src <-chan interface{},
	fn func(arg interface{}) interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-src:
				if v := fn(data); v != nil {
					valueStream <- fn(v)
				}
			}
		}
	}()
	return valueStream
}

func toString(
	ctx context.Context,
	valueStream <-chan interface{},
) <-chan string {
	stringStream := make(chan string)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-ctx.Done():
				return
			case stringStream <- v.(string):
			}
		}
	}()
	return stringStream
}

func toInt(
	ctx context.Context,
	valueStream <-chan interface{},
) <-chan int {
	stringStream := make(chan int)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-ctx.Done():
				return
			case stringStream <- v.(int):
			}
		}
	}()
	return stringStream
}

func fanIn(
	ctx context.Context,
	channels ...<-chan interface{},
) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})
	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-ctx.Done():
				return
			case multiplexedStream <- i:
			}
		}
	}
	// Select from all the channels
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}
	// Wait for all the reads to complete
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()
	return multiplexedStream
}
