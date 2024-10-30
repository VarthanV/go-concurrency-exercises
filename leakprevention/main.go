package main

import (
	"fmt"
	"time"
)

func sampleGoroutineLeak() {
	// doWork: this routine is not garbage collected and lives
	// whole lifetime of the program
	doWork := func(strs <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("do work exited")
			defer close(completed)
			for s := range strs {
				fmt.Println(s)
			}
		}()
		return completed
	}

	doWork(nil)
	fmt.Println("done!!")
}

func mitigatingLeakwithDoneChan() {
	fmt.Println("####################### mitigatingLeakwithDoneChan ###############")
	done := make(chan interface{})

	doWork := func(strs <-chan string, done <-chan interface{}) <-chan interface{} {
		terminate := make(chan interface{})
		go func() {
			defer close(terminate)
			select {
			case <-done:
				fmt.Println("done..............")
				return
			case val := <-strs:
				fmt.Println("val is ", val)
			}
		}()
		return terminate
	}

	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	terminated := doWork(nil, done)
	<-terminated // block until signal
	fmt.Println("################################")
}

func main() {
	sampleGoroutineLeak()
	mitigatingLeakwithDoneChan()
}
