// Package channels
package channels

import (
	"fmt"
	"sync"
)

// Смерджим несколько каналов в один
func mergeChan(cs ...<-chan int) <-chan int {
	out := make(chan int)
	wg := sync.WaitGroup{}

	// горутина для переноса сообщений из одного канала в общий канал out
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	// В отдельной горутине чтобы не блочить возвращение канала out
	go func() { wg.Wait(); close(out) }()
	return out
}

func MergeThreeChannelsToOne() {
	// Создам три канала
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	// Канал out
	o := mergeChan(a, b, c)
	// Канал для ожидания горутины для вывода значений из o
	d := make(chan bool)

	go func(o <-chan int, done chan<- bool) {
		for i := range o {
			fmt.Println(i)
		}
		done <- true
	}(o, d)

	// Примитивно отправляю сообщения в канал и закрываю его.
	a <- 1
	close(a)
	b <- 2
	close(b)
	c <- 3
	close(c)
	<-d
	close(d)
}
