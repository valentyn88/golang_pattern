package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func randStream(done <-chan interface{}, i int) <-chan int {
	randStream := make(chan int)
	go func() {
		defer close(randStream)
		for j:=0;j<i;j++ {
			select {
			case <-done:
				return
			case randStream<-i:
			}
		}
	}()
	return randStream
}

func fanIn(done <-chan interface{}, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	multiplexedStream := make(chan int)

	multiplex := func(c <-chan int) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream<-i:
			}
		}
	}

	//Select from all channels
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	//Wait for all the reads to complete
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

type Taker struct {}

func (Taker) take(done <-chan interface{}, valueStream <-chan int, num int) <-chan int {
	takeStream := make(chan int)
	go func() {
		defer close(takeStream)
		for i:=0;i<num;i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}

func main()  {
	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	numFinders := runtime.NumCPU()
	fmt.Printf("Num finders: %d\n", numFinders)
	finders := make([] <-chan int, numFinders)
	for i:=0;i<numFinders;i++ {
		finders[i] = randStream(done, rand.Intn(100500))
	}
	taker := Taker{}
	for prime := range taker.take(done, fanIn(done, finders...), numFinders) {
		fmt.Printf("Prime %v\n", prime)
	}

	fmt.Printf("Time since: %v\n", time.Since(start))
}
