package main

import (
	"fmt"
	"sync"
)

func Tee(in chan string) (chan string, chan string) {
	out1 := make(chan string)
	out2 := make(chan string)

	go func() {
		defer close(out1)
		defer close(out2)

		for msg := range in {
			out1 <- msg
			out2 <- msg
		}
	}()

	return out1, out2
}

func main() {
	const messageCount int = 5
	in := make(chan string)

	go func() {
		defer close(in)

		for index := range messageCount {
			in <- fmt.Sprintf("Message #%d", index+1)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	{
		chEmail, chSMS := Tee(in)

		go func() {
			defer wg.Done()

			for msg := range chEmail {
				fmt.Printf("send email message: %s\n", msg)
			}
		}()

		go func() {
			defer wg.Done()

			for msg := range chSMS {
				fmt.Printf("send sms message: %s\n", msg)
			}
		}()
	}
	wg.Wait()
}
