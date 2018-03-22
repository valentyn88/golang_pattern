package main

import (
	"fmt"
)

func repeat(done <-chan interface{}, values ...int) <-chan int {
	repeatStream := make(chan int)
	go func() {
		defer close(repeatStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
					case repeatStream <-v:
				}
			}
		}
	}()

	return repeatStream
}

func take(done <-chan interface{},in <-chan int, count int) <-chan int {
	takeStream := make(chan int)
	go func() {
		defer close(takeStream)
		for i:=0;i<count;i++ {
			select {
				case <-done:
					return
					case takeStream <- <-in:
			}
		}
	}()

	return takeStream
}

func tee(done <-chan interface{}, in <-chan int) (<-chan int, <-chan int) {
	out1 := make(chan int)
	out2 := make(chan int)

	go func() {
		defer close(out1)
		defer close(out2)
		for val := range in {
			var out1, out2 = out1, out2
			for i:=0;i<2;i++ {
				select {
				case <-done:
					return
					case out1 <-val:
						out1 = nil
						case out2 <-val:
							out2 = nil
				}
			}
		}
	}()

	return out1, out2
}

func main()  {
	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1,2,3,4), 2))
	for v := range out1 {
		fmt.Printf("out1 val: %d, out2 val: %d\n", v, <-out2)
	}
}
