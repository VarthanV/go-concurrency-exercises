package main

import (
	"fmt"
	"time"
)

func sendingIterationVariablesOnChannel() {
	fmt.Println("############## sendingIterationVariablesOnChannel ##############")
	stream := make(chan string)

	go func() {
		for _, s := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
			stream <- s
		}
	}()

	for {
		select {
		case val := <-stream:
			fmt.Println("val received is ", val)

		case <-time.After(2 * time.Second):
			fmt.Println("Timeout!!")
			fmt.Println("#########################################")
			return
		}
	}

}

func infiniteWaitUntilStopped() {
	fmt.Println("############## sendingIterationVariablesOnChannel ##############")

	done := make(chan bool)

	defer func() {
		fmt.Println("#########################################")
	}()

	go func() {
		time.Sleep(5 * time.Second)
		done <- true
	}()

	for {
		select {
		case <-done:
			return
		default:
			fmt.Println("Stop me folks !!!!!!!")
		}
	}
}
func main() {
	sendingIterationVariablesOnChannel()
	infiniteWaitUntilStopped()
}
