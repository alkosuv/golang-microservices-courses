package main

import (
	"fmt"
	"sync"
)

func Bridge(in chan chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		var wg sync.WaitGroup

		for ch := range in {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for v := range ch {
					out <- v
				}
			}()
		}

		wg.Wait()
	}()

	return out
}

func main() {
	const ChannelSize int = 3
	in := make(chan chan string)

	go func() {
		defer close(in)

		ch1 := make(chan string, ChannelSize)
		defer close(ch1)
		ch2 := make(chan string, ChannelSize)
		defer close(ch2)

		for i := 0; i < ChannelSize; i++ {
			ch1 <- fmt.Sprintf("ch1 message: #%d", i)
			ch2 <- fmt.Sprintf("ch2 message: #%d", i)
		}

		in <- ch1
		in <- ch2
	}()

	out := Bridge(in)
	for v := range out {
		fmt.Println(v)
	}
}
