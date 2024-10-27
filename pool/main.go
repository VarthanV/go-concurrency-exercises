package main

import (
	"fmt"
	"sync"
)

func basicExample() {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("creating new instace")
			return struct{}{}
		},
	}

	myPool.Get() // Invokes new callback
	instance := myPool.Get()
	myPool.Put(instance) // Put back the resource acquire increase no of resource
	// acquired to 1
	myPool.Get() // reuse instance already allocated
}

func realWorldComplexUseCase() {
	var (
		numInstancesCreated int
		wg                  sync.WaitGroup
	)
	calcPool := &sync.Pool{
		New: func() interface{} {
			numInstancesCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	// Seed the pool with 4kb

	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	numWorkers := 1024 * 1024
	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)
			// asumme something interesting done
		}()
	}
	wg.Wait()
	fmt.Printf("%d calculators were created\n", numInstancesCreated) // results are non
	// deterministic depends on CPU cycle.
}
func main() {
	basicExample()
	realWorldComplexUseCase()
}
