// Шаблон FanIn (Мультиплексор) в Go используется для объединения результатов из нескольких источников (горутин) в один канал.
// Это позволяет эффективно собирать данные, которые обрабатываются параллельно,
// и отправлять их в одно место для дальнейшей обработки или вывода.
package main

import (
	"fmt"
	"sync"
	"time"
)

func FanIn(in ...<-chan string) <-chan string {
	out := make(chan string)

	var wg sync.WaitGroup
	for _, ch := range in {
		wg.Add(1)

		go func(ch <-chan string) {
			defer wg.Done()

			for v := range ch {
				out <- fmt.Sprintf("pressed: %s", v)
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	const (
		ChannelSize  int = 3
		MessageCount int = 5
	)

	in := make([]<-chan string, ChannelSize)
	for i := 0; i < ChannelSize; i++ {
		ch := make(chan string)
		in[i] = ch

		go func(ch chan string, indexCh int) {
			defer close(ch)

			for indexMessage := range MessageCount {
				ch <- fmt.Sprintf("channel ID: %d message ID: %d", indexCh, indexMessage)
				time.Sleep(time.Second)
			}
		}(ch, i)
	}

	out := FanIn(in...)
	for message := range out {
		fmt.Println(message)
	}
}
