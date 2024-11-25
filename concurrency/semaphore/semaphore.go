package main

import (
	"fmt"
	"sync"
	"time"
)

// Semaphore представляет собой примитив синхронизации, который управляет доступом к ограниченным ресурсам,
// позволяя контролировать количество одновременно работающих горутин.

type Semaphore struct {
	ch chan struct{}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func main() {
	const count int = 2
	wg := sync.WaitGroup{}

	s := &Semaphore{make(chan struct{}, count)}

	for i := 0; i < count+1; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			fmt.Printf("wiat index: %d\n", index)

			s.Acquire()
			defer s.Release()

			fmt.Printf("start index: %d\n", index)
			time.Sleep(time.Second * time.Duration(10))
			fmt.Printf("stop index: %d\n", index)
		}(i)
	}

	wg.Wait()
}
