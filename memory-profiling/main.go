package main

import (
	"fmt"
	"runtime"
	"sync"
)

/*
A newly minted goroutine is given a few kilobytes,
which is almost always enough. When it isn’t,
the run-time grows (and shrinks) the memory for storing the stack automatically, allowing many goroutines to live in a modest amount of memory.
*/
func main() {

	var wg sync.WaitGroup
	var c <-chan interface{}

	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	noop := func() {
		wg.Done()
		<-c
	}

	const numGoroutines = 1e4
	wg.Add(numGoroutines)

	before := memConsumed()

	for i := 0; i < numGoroutines; i++ {
		go noop()
	}

	wg.Wait()

	after := memConsumed()

	fmt.Printf("%.3fkb\n", float64(after-before)/numGoroutines/1000)
}
