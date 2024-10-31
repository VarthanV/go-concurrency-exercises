package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id" gorm:"primaryKey"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func generator(ctx context.Context, urls ...string) <-chan string {
	urlStream := make(chan string)
	go func() {
		defer close(urlStream)
		for _, url := range urls {
			select {
			case <-ctx.Done():
				return
			case urlStream <- url:
			}
		}
	}()
	return urlStream
}

func doHTTP(ctx context.Context, urlStream <-chan string) <-chan *Todo {
	resultStream := make(chan *Todo)
	go func() {
		defer close(resultStream)
		for url := range urlStream {
			select {
			case <-ctx.Done():
				return
			default:
				resp, err := http.Get(url)
				if err != nil {
					log.Printf("Failed to fetch URL %s: %v", url, err)
					continue
				}

				var todo Todo
				if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
					log.Printf("Failed to decode JSON for URL %s: %v", url, err)
					continue
				}
				resultStream <- &todo
				resp.Body.Close()
			}
		}
	}()
	return resultStream
}

func insertInDB(ctx context.Context, todos <-chan *Todo, db *gorm.DB, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	for todo := range todos {
		select {
		case <-ctx.Done():
			return
		default:
			err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&todo).Error
			if err != nil {
				errChan <- fmt.Errorf("DB insertion error for ID %d: %w", todo.ID, err)
			}
		}
	}
}

func WebScrapperPipelineDriver() {
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&Todo{}); err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	urlStream := generator(ctx,
		"https://jsonplaceholder.typicode.com/posts/1",
		"https://jsonplaceholder.typicode.com/posts/2",
		"https://jsonplaceholder.typicode.com/posts/3",
	)

	errChan := make(chan error)
	var wg sync.WaitGroup

	todoStream := doHTTP(ctx, urlStream)
	wg.Add(1)
	go insertInDB(ctx, todoStream, db, &wg, errChan)

	go func() {
		for err := range errChan {
			log.Println("Error:", err)
		}
	}()

	wg.Wait()
	close(errChan)
}
