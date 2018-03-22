package main

import "fmt"

func main()  {
	done := make(chan interface{})
	defer close(done)

	data := []int{1,2,3,4,5,6}
	streamFunc := func(done <-chan interface{}, i ...int) <-chan int {
		c := make(chan int)
		go func() {
			defer close(c)
			for _, v := range i {
				select {
				case <-done:
					return
				case c<-v:
				}
			}
		}()

		return c
	}

	orDone := func(done <-chan interface{}, c <-chan int) <-chan int {
		valStream := make(chan int)
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}

				}
			}
		}()

		return valStream
	}

	for v := range orDone(done, streamFunc(done, data...)) {
		fmt.Println(v)
	}
}
