package main

import (
	"fmt"
	"sync"
	"time"
)

type Barrier struct {
	count     int        // количество горутин
	waiting   int        // текущее количество горутин, которые находят в ожидании
	mx        sync.Mutex // mutex для синхронизации доступа
	mxRefresh sync.Mutex // mutex для заврешения ожидания
	cond      *sync.Cond // переменная для ожидания

	waitingGoroutineIDs []int
}

func NewBarrier(count int) *Barrier {
	b := &Barrier{
		count:               count,
		waitingGoroutineIDs: make([]int, 0, count),
	}
	b.cond = sync.NewCond(&b.mx)

	go func() {
		t := time.Second * time.Duration(5)

		ticker := time.NewTicker(t)
		for range ticker.C {
			if b.waiting == 0 {
				ticker.Reset(t)
				continue
			}

			b.refresh()
			ticker.Reset(t)
		}
	}()

	return b
}

func (b *Barrier) refresh() {
	b.mxRefresh.Lock()

	fmt.Printf("waiting goroutines: %v\n", b.waitingGoroutineIDs)

	b.waiting = 0
	b.waitingGoroutineIDs = make([]int, 0, b.count)
	b.cond.Broadcast()

	b.mxRefresh.Unlock()
}

func (b *Barrier) Await(goroutineID int) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.waitingGoroutineIDs = append(b.waitingGoroutineIDs, goroutineID)
	b.waiting++

	if b.waiting == b.count {
		b.refresh()
		return
	}

	b.cond.Wait()
}

func main() {
	const goroutinesCount int = 3
	barrier := NewBarrier(goroutinesCount)

	var wg sync.WaitGroup
	for index := 0; index < 10; index++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			fmt.Printf("Goroutine #%d doing something...\n", index)
			barrier.Await(index)
			fmt.Printf("Goroutine #%d done waiting for other goroutines\n", index)
		}()
	}

	wg.Wait()
}
