package main

import "context"

func bridge(
	ctx context.Context,
	chanStream <-chan<-chan interface{},
) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-ctx.Done():
				return
			}
			for val := range orDone(ctx, stream) {
				valStream <- val
			}
		}
	}()
	return valStream
}
