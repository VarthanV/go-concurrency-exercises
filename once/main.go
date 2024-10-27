package main

import (
	"fmt"
	"sync"
)

func basicExample() {
	var (
		wg    sync.WaitGroup
		once  sync.Once
		count = 0
	)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(func() {
				count++
			}) // this is executed excatly once
		}()
	}

	wg.Wait()
	fmt.Println("count is ", count)
}

func main() {
	basicExample()
}
