// Паттерн FanOut (Демультиплексор) в Go используется для распределения работы между несколькими горутинами,
// что позволяет параллельно обрабатывать задачи и увеличивать пропускную способность приложения.
//
// Основные этапы работы паттерна fan-out
//  1. Разделение задачи: Основная задача делится на несколько подзадач, которые могут выполняться параллельно.
//     Каждая подзадача назначается отдельной горутине.
//  2. Параллельное выполнение: Горутину, отвечающую за выполнение подзадачи, можно запустить для обработки данных из общего канала.
//  3. Увеличение производительности: Параллельная обработка позволяет значительно сократить время выполнения задач, особенно если они ресурсоемкие.
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
