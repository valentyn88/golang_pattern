package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exit...")
			defer close(terminated)
			for {
				select {
				case <-done:
					return
				case <-strings:
					s := <-strings
					fmt.Printf("printing strings chanel... %s", s)
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("close done chanel...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done.")
}
