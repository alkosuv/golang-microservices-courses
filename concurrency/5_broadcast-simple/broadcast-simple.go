// Паттерн Broadcast в Go используется для отправки сообщений нескольким подписчикам одновременно.
// Это позволяет одной горутине или источнику данных рассылать информацию сразу в несколько мест,
// что удобно для реализации систем уведомлений, событий и других сценариев, требующих одновременной доставки данных.
//
// Основные принципы работы паттерна broadcast
// 1. Множественные подписчики: Один источник данных может иметь несколько подписчиков, каждый из которых получает копию сообщения.
// 2. Отправка сообщений: Когда источник отправляет сообщение, оно дублируется и отправляется всем подписчикам.
// 3. Использование каналов: В Go этот паттерн часто реализуется с помощью каналов и горутин.
package main

import (
	"fmt"
	"sync"
)

func notifier(signals chan int) {
	close(signals)
}

func subscriber(signals chan int) {
	<-signals
	fmt.Println("signaled")
}

func main() {
	signals := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		subscriber(signals)
	}()

	go func() {
		defer wg.Done()
		subscriber(signals)
	}()

	go func() {
		defer wg.Done()
		notifier(signals)
	}()

	wg.Wait()
}
