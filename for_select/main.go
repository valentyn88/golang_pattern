package main

import (
	"fmt"
	"time"
)

func main() {
	printData := func(done <-chan bool, data []string) {
		for {
			select {
			case <-done:
				return
			default:
			}

			for _, v := range data {
				fmt.Println(v)
			}
		}
	}

	stopFunc := func(done chan<- bool) {
		go func() {
			time.Sleep(time.Second * 1)
			done <- true
		}()
	}

	data := []string{"a", "b", "c", "d"}
	done := make(chan bool)

	stopFunc(done)
	printData(done, data)
}
