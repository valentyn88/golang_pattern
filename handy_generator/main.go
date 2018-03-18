package main

import (
	"fmt"
	"math/rand"
)

type Repeater struct {}

func (Repeater) repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	repeatStream := make(chan interface{})
	go func() {
		defer close(repeatStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
					case repeatStream <- v:
				}
			}
		}
	}()

	return repeatStream
}

type Taker struct {}

func (Taker) take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})
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

type RepeateFn struct {}

func (RepeateFn) repeat(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	repeatStream := make(chan interface{})
	go func() {
		defer close(repeatStream)
		select {
		case <-done:
			return
			case repeatStream<-fn():
		}
	}()

	return repeatStream
}

func main()  {
	done := make(chan interface{})
	defer close(done)

	t := Taker{}
	r := Repeater{}
	rFn := RepeateFn{}
	fn := func() interface{} {return rand.Int()}

	for v := range t.take(done, r.repeat(done, 1), 10) {
		fmt.Println(v)
	}

	for v := range t.take(done, rFn.repeat(done, fn), 5) {
		fmt.Println(v)
	}
}
