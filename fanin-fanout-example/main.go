// generator() -> square() -> print

package main

import (
	"fmt"
	"sync"
)

func generator(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	// Implement fan-in
	// merge a list of channels to a single channel
	out := make(chan int)
	var wg sync.WaitGroup
	
	getChannelNums := func(ch <-chan int) {
		for n := range ch {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, ch := range cs {
		go getChannelNums(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	in := generator(2, 3, 4, 5)

	// fan out square stage to run two instances.
	sq1 := square(in)
	sq2 := square(in)

	// fan in the results of square stages.
	for n := range merge(sq1, sq2) {
		fmt.Println(n)
	}
}