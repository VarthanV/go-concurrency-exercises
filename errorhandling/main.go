package main

import (
	"fmt"
	"log"
	"net/http"
)

type Result struct {
	Error    error
	Response *http.Response
}

func main() {

	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Response: resp, Error: err}

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

	urls := []string{"https://google.com", "https://duckduckgo.com", "https://bas"}

	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			log.Println("error in checking status ", result.Error.Error())
			continue
		}

		fmt.Println("Response status code ", result.Response.StatusCode)

	}
}
