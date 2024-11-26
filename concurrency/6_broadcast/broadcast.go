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
	"time"
)

type Broadcast struct {
	source           <-chan string
	subscribers      []chan string
	addSubscriber    chan chan string
	removeSubscriber chan (<-chan string)
}

func NewBroadcast(source <-chan string) *Broadcast {
	b := &Broadcast{
		source:           source,
		subscribers:      make([]chan string, 0),
		addSubscriber:    make(chan chan string),
		removeSubscriber: make(chan (<-chan string)),
	}

	go b.run()
	return b
}

func (b *Broadcast) run() {
	for {
		select {
		case s := <-b.addSubscriber:
			b.subscribers = append(b.subscribers, s)

		case subscriber := <-b.removeSubscriber:
			for i, s := range b.subscribers {
				if s == subscriber {
					close(s)
					b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
					break
				}
			}

		case msg, ok := <-b.source:
			if !ok {
				for _, subscriber := range b.subscribers {
					close(subscriber)
				}
				return
			}

			for _, subscriber := range b.subscribers {
				subscriber <- msg
			}
		}
	}
}

func (b *Broadcast) AddSubscribe() <-chan string {
	subscribe := make(chan string)
	b.addSubscriber <- subscribe
	return subscribe
}

func (b *Broadcast) RemoveSubscribe(subscriber <-chan string) {
	b.removeSubscriber <- subscriber
}

func main() {
	const (
		subscriberSize int = 3
		MessageCount   int = 10
	)

	source := make(chan string)
	broadcast := NewBroadcast(source)

	go func() {
		defer close(source)
		for i := 0; i < MessageCount; i++ {
			source <- fmt.Sprintf("Message %d", i)
			time.Sleep(time.Second)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(subscriberSize)

	for index := 0; index < subscriberSize; index++ {
		subscriber := broadcast.AddSubscribe()

		go func(subscriber <-chan string, indexCh int) {
			defer wg.Done()

			for msg := range subscriber {
				fmt.Printf("Subscriber %d received message %s\n", indexCh, msg)
			}
		}(subscriber, index)

		go func(indexCh int) {
			if indexCh%2 == 0 {
				time.Sleep(time.Second * time.Duration(3))
				broadcast.RemoveSubscribe(subscriber)
			}
		}(index)

	}

	wg.Wait()
}
