package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func work(id int, wg *sync.WaitGroup, inputStream <-chan int, outStream chan<- int) {
	go func() {
		defer wg.Done()
		log.Println("Spawned worker id ", id)
		for val := range inputStream {
			time.Sleep(2 * time.Second) // simulate work
			outStream <- val * 2
		}
	}()
}

func queueing() {
	var (
		numWorkers = 2
		numJobs    = 10
		wg         sync.WaitGroup
	)

	jobStream := make(chan int, numWorkers)
	outStream := make(chan int, numWorkers)

	// Spawn workers

	for i := 0; i < numWorkers; i++ {
		// Spawn workers
		wg.Add(1)
		work(i, &wg, jobStream, outStream)
	}

	// Feed jobs
	go func() {
		defer close(jobStream)

		for i := 0; i < numJobs; i++ {
			jobStream <- i
		}
	}()

	// Wait until all work is done and close result chan
	go func() {
		wg.Wait()
		close(outStream)
	}()

	// Range  on result  stream
	for out := range outStream {
		fmt.Println("Received out is ", out)
	}
}

func main() {
	queueing()
}
