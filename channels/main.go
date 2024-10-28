package main

import (
	"fmt"
	"sync"
)

func basicReadWriteExample() {
	fmt.Println("############## basicReadWriteExample ##############")
	stringStream := make(chan string)
	go func() {
		stringStream <- "Hello world" // Sending data to stream
	}()

	fmt.Println(<-stringStream) // receiving data and waiting for it
	fmt.Println("##############################")
}

func channelReturnsBoolValWhenReading() {
	fmt.Println("############## channelReturnsBoolValWhenReading ##############")
	stringStream := make(chan string)
	go func() {
		stringStream <- "Hello world" // Sending data to stream
	}()

	val, ok := <-stringStream
	fmt.Printf("(%v): %v\n", ok, val)
	fmt.Println("##############################")
}

func readingFromClosedChan() {
	fmt.Println("############## readingFromClosedChan ##############")

	intStream := make(chan int)
	close(intStream)

	val, ok := <-intStream
	fmt.Printf("(%v): %v\n", ok, val) // falsy value of int will be returned , ok will be false
	fmt.Println("##############################")

}

func rangingOverChannel() {
	intStream := make(chan int)

	fmt.Println("############## rangingOverChannel ##############")

	go func() {
		defer close(intStream) // closing in defer is common pattern
		for i := 0; i < 10; i++ {
			intStream <- i
		}
	}()

	for val := range intStream {
		fmt.Println(val)
	}

	fmt.Println("##############################")
}

func signallingMechanism() {
	fmt.Println("############## signallingMechanism ##############")
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin // wait until unblocked
			fmt.Println(i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Wait()

	fmt.Println("##############################")

}

func main() {
	basicReadWriteExample()
	channelReturnsBoolValWhenReading()
	readingFromClosedChan()
	rangingOverChannel()
	signallingMechanism()
}
