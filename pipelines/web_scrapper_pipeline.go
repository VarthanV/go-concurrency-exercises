package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id" gorm:"primaryKey" `
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

		for url := range urlStream {
			select {
			default:
				var (
					httpResp *Todo
				)

				log.Println("Fetching url ", url)
				resp, err := http.Get(url)
				if err != nil {
					resultStream <- Process{
						Err: err,
					}
					continue
				}

				// Read request body
				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					resultStream <- Process{
						Err: err,
					}
					continue
				}

				err = json.Unmarshal(respBody, &httpResp)
				if err != nil {
					resultStream <- Process{
						Err: err,
					}
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

func insertInDB(done <-chan interface{}, processStream <-chan Process, db *gorm.DB) <-chan Process {
	resultStream := make(chan Process)

	go func() {
		defer close(resultStream)

		for val := range processStream {
			select {
			default:
				if val.Err != nil {
					resultStream <- val
					continue
				}

				if val.Todo != nil {
					log.Println("inserting into db with id ", val.Todo.ID)
					err := db.Model(&Todo{}).
						Clauses(clause.OnConflict{DoNothing: true}).
						Create(&val.Todo).Error
					if err != nil {
						val.Err = errors.Join(val.Err, err)
						resultStream <- val
						continue
					}
					resultStream <- val
				}

			case <-done:
				return
			}
		}
	}()

	return resultStream
}

func WebScrapperPipelineDriver() {
	var (
		errChan = make(chan error, 5)
		wg      sync.WaitGroup
	)

	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("unable to open db ", err)
	}

	err = db.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatal("unable to automigrate ", err)
	}

	done := make(chan interface{})
	defer close(done)

	urlStream := generator(done,
		"https://jsonplaceholder.typicode.com/posts/1",
		"https://jsonplaceholder.typicode.com/posts/2",
		"https://jsonplaceholder.typicode.com/posts/3",
		"https://bas",
	)

	logErrorToFile := func(done <-chan interface{}, errChan <-chan error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for err := range errChan {
				select {
				default:
					log.Println(err)

				case <-done:
					return
				}
				// TODO: will log in file later
			}

		}()
	}

	logErrorToFile(done, errChan)

	pipeline := insertInDB(done, doHTTP(done, urlStream), db)

	for val := range pipeline {
		if val.Err != nil {
			errChan <- val.Err
		}
	}
	close(errChan)
	wg.Wait()

}
