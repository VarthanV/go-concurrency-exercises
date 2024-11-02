package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/fatih/color"
)

type Todo struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Result struct {
	WorkerID int
	Todo     Todo
	Error    error
}

func fanOut() {
	var wg sync.WaitGroup

	makeRequest := func(id int, ctx context.Context, wg *sync.WaitGroup, inputStream <-chan string, outStream chan<- Result) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.Println("quitting worker id ", id)
				return
			case url, ok := <-inputStream:

				if !ok {
					return
				}

				log.Println("executing url ", url)
				httpResult := Todo{}
				result := Result{
					WorkerID: id,
				}

				req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // Use GET method
				if err != nil {
					log.Println("error in making request ", err)
					result.Error = err
					outStream <- result
					continue
				}

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Println("error in doing request ", err)
					result.Error = err
					outStream <- result
					continue
				}
				defer res.Body.Close() // Ensure response body is closed

				respBody, err := io.ReadAll(res.Body)
				if err != nil {
					log.Println("error in reading resp body ", err)
					result.Error = err
					outStream <- result
					continue
				}

				err = json.Unmarshal(respBody, &httpResult)
				if err != nil {
					log.Println("error in unmarshalling ", err)
					result.Error = err
					outStream <- result
					continue
				}

				result.Todo = httpResult
				outStream <- result
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	urlChan := make(chan string)
	resultChan := make(chan Result)

	// Start workers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go makeRequest(i, ctx, &wg, urlChan, resultChan)
	}

	go func() {
		for i := 1; i <= 10; i++ {
			urlChan <- fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", i)
		}
		close(urlChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for val := range resultChan {
		if val.Error != nil {
			color.Red("Worker %d encountered error: %v\n", val.WorkerID, val.Error)
		} else {
			color.Green("Worker %d retrieved Todo: %+v\n", val.WorkerID, val.Todo)
		}
	}

}

func main() {
	fanOut()
}
