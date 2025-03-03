package main

import "fmt"

func main() {
	// Here we instantiate the channel within the lexical scope of the chanOwner function.
	// This limits the scope of the write aspect of the results channel to the closure
	// defined below it. In other words, it confines the write aspect of this channel
	// to prevent other goroutines from writing to it.
	chanOwner := func() <-chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 20; i++ {
				results <- i
			}
		}()
		return results
	}

	// 	Here we receive the read aspect of the channel and we’re able to pass it into the
	// consumer, which can do nothing but read from it. Once again this confines the
	// main goroutine to a read-only view of the channel.
	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}
	results := chanOwner()
	consumer(results)
}
