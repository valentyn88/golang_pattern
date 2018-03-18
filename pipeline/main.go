package main

import (
	"fmt"
)

type Pipeline interface {
	do(done <-chan interface{}, intStream <-chan int, param int) <-chan int
}

type Multiplier struct {}

func (Multiplier) do(done <-chan interface{}, intStream <-chan int, param int) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		defer close(multipliedStream)
		for i := range intStream {
			select {
			case <-done:
				return
				case multipliedStream <- i * param:
			}
		}
	}()

	return multipliedStream
}

type Adder struct {}

func (Adder) do(done <-chan interface{}, intStream <-chan int, param int) <-chan int {
	addStream := make(chan int)
	go func() {
		defer close(addStream)
		for i := range intStream {
			select {
			case <-done:
				return
				case addStream <- i + param:
			}
		}
	}()

	return addStream
}

type Generator struct {}

func (Generator) generate (done <-chan interface{},ints ...int) <-chan int {
	intStream := make(chan int, len(ints))
	go func() {
		defer close(intStream)
		for _, i := range ints {
			select {
			case <-done:
				return
			case intStream <- i:
			}
		}
	}()

	return intStream
}

func main()  {
	done := make(chan interface{})
	defer close(done)

	gen := Generator{}
	mul := Multiplier{}
	add := Adder{}
	intStream := gen.generate(done, 1,2,3,4)
	pipeline := mul.do(done, add.do(done,mul.do(done, intStream, 2), 1),2)
	for i := range pipeline {
		fmt.Println(i)
	}
}
