package main

import (
	"context"
	"testing"
)

func BenchmarkGeneric(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b.ResetTimer()
	for range toString(ctx, take(ctx, repeat(ctx, "a"), b.N)) {
	}
}

func BenchmarkTyped(b *testing.B) {
	repeat := func(done <-chan interface{}, values ...string) <-chan string {
		valueStream := make(chan string)
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	take := func(
		done <-chan interface{},
		valueStream <-chan string,
		num int,
	) <-chan string {
		takeStream := make(chan string)
		go func() {
			defer close(takeStream)
			for i := num; i > 0 || i == -1; {
				if i != -1 {
					i--
				}
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}
	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()
	for range take(done, repeat(done, "a"), b.N) {
	}
}
