// Semaphore представляет собой примитив синхронизации, который управляет доступом к ограниченным ресурсам,
// позволяя контролировать количество одновременно работающих горутин.
package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	ch chan struct{}
}

// Acquire метод для полученеия доступа к ресурсам
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Release метод для завершения потребления ресурсов
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
