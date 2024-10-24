package main

import (
	"fmt"
	"sync"
)

func main() {

	basicIncrementDecrementMuxtexExample := func() {
		var count int

		var (
			mu sync.Mutex
			wg sync.WaitGroup
		)

		// Increment
		increment := func() {
			mu.Lock()
			defer mu.Unlock()
			count++
			fmt.Println("Incrementing count ", count)
		}

		// Decrement
		decrement := func() {
			mu.Lock()
			defer mu.Unlock()
			count--
			fmt.Println("Decrementing count ", count)
		}

		// Call increment n times
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				increment()
			}()
		}
		// Call decrement n times

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				decrement()
			}()
		}

		wg.Wait()
	}

	basicIncrementDecrementMuxtexExample()
}
