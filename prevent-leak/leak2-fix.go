package main

import (
	"fmt"
	"math/rand"
)

func main() {
	newRandStream := func(done <-chan struct{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer close(randStream)
			defer fmt.Println("newRandStream closure exited.")
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					fmt.Println("going to break")
					return
				}
			}
		}()
		return randStream
	}

	done := make(chan struct{})
	randStream := newRandStream(done)

	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	done <- struct{}{}
	<-randStream
}
