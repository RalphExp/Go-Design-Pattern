package main

import (
	"context"
)

func generator(ctx context.Context, integers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for _, i := range integers {
			select {
			case <-ctx.Done():
				return
			case intStream <- i:
			}
		}
	}()
	return intStream
}

func multiply(ctx context.Context, intStream <-chan int, multiplier int) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		defer close(multipliedStream)
		for i := range intStream {
			select {
			case <-ctx.Done():
				return
			case multipliedStream <- i * multiplier:
			}
		}
	}()
	return multipliedStream
}

func add(ctx context.Context, intStream <-chan int, additive int) <-chan int {
	addedStream := make(chan int)
	go func() {
		defer close(addedStream)
		for i := range intStream {
			select {
			case <-ctx.Done():
				return
			case addedStream <- i + additive:
			}
		}
	}()
	return addedStream
}
