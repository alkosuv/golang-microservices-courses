// Паттерн Tee в Go, аналогично команде `tee` в Unix, используется для дублирования данных из одного канала в несколько
// выходных каналов. Это позволяет одной горутине отправлять данные в несколько мест одновременно,
// что может быть полезно для логирования, мониторинга или параллельной обработки данных.
//
// Основные принципы работы паттерна tee
//  1. Дублирование данных: Паттерн tee принимает данные из одного канала и отправляет их в два или более выходных каналов.
//  2. Независимость: Каждое из выходных мест может обрабатывать данные независимо, что позволяет избежать блокировок
//     и задержек при отправке данных.
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
