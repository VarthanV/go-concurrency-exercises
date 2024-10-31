# Error handling

- In concurrent programs , error handling can be difficult to get right

- The question arises  , **Who should be responsible for handling the error?**

- In general concurrent process should send their errors to another part of program that has complete information about the state of the program.

## Example

```go
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

```

- We need to seperate concerns of error handling from our producer goroutine.

- The goroutine that spawned the producer goroutine has more context about the running program and can make ``intelligent desicions`` on what to do with the error.

- Errors should be considered first class citizens when constructing values to return from goroutine.

- If goroutines can produce errors, they should be tightly coupled with result type and should be passed in return similar to how we do in synchornous functions.
