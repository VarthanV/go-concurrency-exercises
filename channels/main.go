package main

import (
	"fmt"
	"sync"
	"time"
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

func bufferedChan() {
	fmt.Println("############## bufferedChan ##############")
	intStream := make(chan int, 4)

	go func() {
		defer close(intStream)
		for i := 0; i < 5; i++ {
			intStream <- i
		}
	}()

	for val := range intStream {
		fmt.Println(val)
	}

	fmt.Println("##############################")

}

// producerConsumerPattern: producer owns the responsiblity and provides a ready only stream of channel
// for the consumer to consume
func producerConsumerPattern() {
	fmt.Println("############## producerConsumerPattern ##############")

	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5)
		go func() {
			defer close(resultStream)
			for i := 0; i < 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream
	}

	resultStream := chanOwner()

	for result := range resultStream { // will block here
		fmt.Println("received ", result)
	}

	fmt.Println("done receiving!!")
	fmt.Println("##############################")

}

func basicSelect() {
	fmt.Println("############## basicSelect ##############")

	c1 := make(chan interface{})
	c2 := make(chan interface{})

	var c3 chan<- interface{}

	go func() {
		time.Sleep(3 * time.Second)
		c2 <- "foo"
	}()

	select {
	case <-c1:
		fmt.Println("received from c1")
	case <-c2:
		fmt.Println("received from c2")
	case c3 <- true:
		// do something
	}

	fmt.Println("exiting!!")

}

func timeOutToPreventBlocking() {
	var c <-chan int
	fmt.Println("############## timeOutToPreventBlocking ##############")

	select {
	case <-c:
		// do something
	case <-time.After(2 * time.Second):
		fmt.Println("timed out!!")
	}
	fmt.Println("##############################")

}

func foreverBlocking() {
	/*
			fatal error: all goroutines are asleep - deadlock!

		goroutine 1 [select (no cases)]:
		main.foreverBlocking(...)
		main.main()
		exit status 2
	*/
	select {}
}

func main() {
	basicReadWriteExample()
	channelReturnsBoolValWhenReading()
	readingFromClosedChan()
	rangingOverChannel()
	signallingMechanism()
	bufferedChan()
	producerConsumerPattern()
	basicSelect()
	timeOutToPreventBlocking()
	foreverBlocking()

}
