package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	err      error
	response *http.Response
}

func main() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)

		go func() {
			defer close(results)

			for _, url := range urls {
				resp, err := http.Get(url)
				result := Result{err: err, response: resp}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()

		return results
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"http://google.com", "badhost.com"}
	for result := range checkStatus(done, urls...) {
		if result.err != nil {
			fmt.Printf("Error occured %v\n", result.err)
			continue
		}
		fmt.Printf("Status code: %v\n", result.response.StatusCode)
	}
}
