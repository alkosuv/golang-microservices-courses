// Шаблон Fan-Out (Демультиплексор) распределяет сообщения из
// одного входного канала между несколькими выходными каналами

package main

import (
	"fmt"
	"sync"
)

func FunOut(in <-chan string, countOutChannel int) []<-chan string {
	out := make([]<-chan string, countOutChannel)

	for index := 0; index < countOutChannel; index++ {
		ch := make(chan string)
		out[index] = ch

		go func(ch chan string, indexCh int) {
			defer close(ch)

			for v := range in {
				ch <- fmt.Sprintf("channel ID %d: message: %s", indexCh, v)
			}
		}(ch, index)
	}

	return out
}

func main() {
	const (
		MessageCount int = 30
		ChannelSize  int = 3
	)

	in := make(chan string)
	out := FunOut(in, ChannelSize)

	go func() {
		for index := 0; index < MessageCount; index++ {
			in <- fmt.Sprintf("message %d", index)
		}

		close(in)
	}()

	var wg sync.WaitGroup
	wg.Add(len(out))

	for _, ch := range out {
		go func() {
			defer wg.Done()

			for v := range ch {
				fmt.Println(v)
			}
		}()
	}

	wg.Wait()
}
