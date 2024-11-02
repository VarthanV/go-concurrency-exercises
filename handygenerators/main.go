package main

import (
	"fmt"
	"math/rand"
	"time"
)

// repeat: Repeats values infinitely until we tell to stop
func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	outStream := make(chan interface{})
	go func() {
		defer close(outStream)
		for {
			for _, val := range values {
				select {
				case <-done:
					return
				case outStream <- val:

				}
			}
		}
	}()

	return outStream
}

// repeatDriver: Driver code to run the repeat example
func repeatDriver() {
	fmt.Println("############### repeatDriver ###################")
	done := make(chan interface{})

	stream := repeat(done, 1, 2, 3, 4, 5)
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	for val := range stream {
		fmt.Println("val is ", val)
	}

	fmt.Println("done streaming!!!")
	fmt.Println("####################################")
}

// take: Takes first nums values from the stream and processes something with it and returns it back
// to the outStream
func take(done <-chan interface{}, stream <-chan interface{}, nums int) <-chan interface{} {
	outStream := make(chan interface{})
	go func() {
		defer close(outStream)
		for i := 0; i < nums; i++ {
			select {
			case <-done:
				return

			case outStream <- <-stream: // Take value from stream and send it immmediately to outstream

			}
		}
	}()

	return outStream
}

func takeDoneDriver() {
	fmt.Println("################### takeDoneDriver #################### ")
	done := make(chan interface{})
	defer close(done)

	for val := range take(done, repeat(done, 1), 5) {
		fmt.Println(val)
	}
	fmt.Println("####################################")
}

// repeatFn: repeats a fn until stopped and passes down the result to a outputstream
func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	outStream := make(chan interface{})
	go func() {
		defer close(outStream)
		for {
			select {
			case <-done:
				return
			case outStream <- fn():
			}
		}
	}()
	return outStream
}

func repateFnDriver() {
	fmt.Println("################### repateFnDriver #################### ")

	done := make(chan interface{})
	defer close(done)

	repeatRandInt := func() interface{} {
		return rand.Int()
	}

	for val := range take(done, repeatFn(done, repeatRandInt), 11) {
		fmt.Println(val)
	}
	fmt.Println("####################################")
}

func fnWithTypeAssertionStage() {
	toInt := func(done, inputStream <-chan interface{}) <-chan int {
		intStream := make(chan int)

		go func() {
			defer close(intStream)
			for val := range inputStream {
				select {
				case <-done:
					return
				case intStream <- val.(int):
				}
			}

		}()
		return intStream
	}

	fmt.Println("################### fnWithTypeAssertionStage #################### ")

	done := make(chan interface{})
	defer close(done)

	for val := range toInt(done, take(done, repeat(done, 1, 2, 3), 3)) {
		fmt.Println(val * 2)
	}

	// outputs 2,4,6
	fmt.Println("####################################")

}

func main() {
	repeatDriver()
	takeDoneDriver()
	repateFnDriver()
	fnWithTypeAssertionStage()
}
