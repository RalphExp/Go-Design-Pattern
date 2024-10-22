package main

import "context"

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
