package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Process struct {
	Todo *Todo
	Err  error
}

func generator(done <-chan interface{}, urls ...string) <-chan string {
	urlStream := make(chan string)

	go func() {
		defer close(urlStream)
		for _, val := range urls {
			select {
			case urlStream <- val:
			case <-done:
				return
			}
		}
	}()
	return urlStream
}

// Stage 1
func doHTTP(done <-chan interface{}, urlStream <-chan string) <-chan Process {
	resultStream := make(chan Process)

	go func() {
		defer close(resultStream)

		for {
			select {
			case url := <-urlStream:
				var (
					httpResp *Todo
				)
				resp, err := http.Get(url)
				if err != nil {
					log.Println("error in doing get request ", err)
					resultStream <- Process{
						Err: err,
					}
					continue
				}

				// Read request body
				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println("error in reading resp body ", err)
					resultStream <- Process{
						Err: err,
					}
					continue
				}

				err = json.Unmarshal(respBody, &httpResp)
				if err != nil {
					log.Println("error in unmarshalling body ", err)
					resultStream <- Process{
						Err: err,
					}
					continue
				}

				resultStream <- Process{
					Todo: httpResp,
				}

			case <-done:
				return
			}
		}
	}()

	return resultStream
}

// Stage 2 insert in db

func insertInDB(done <-chan interface{}, processStream <-chan Process) <-chan Process {
	resultStream := make(chan Process)

	go func() {
		defer close(resultStream)

		for {
			select {
			case val := <-processStream:
				if val.Err != nil {
					log.Println("cant go further in pipeline ", val.Err)
					// Can do meaningful stuffs if required to handle error
					continue
				}

				if val.Todo != nil {
					// TODO: insert in db
				}

			case <-done:
				return
			}
		}
	}()

	return resultStream
}
