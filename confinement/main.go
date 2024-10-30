package main

import "fmt"

func adhocConfinement() {
	fmt.Println("############## adhocConfinement ##############")

	data := make(chan int, 4)

	loopData := func(ch chan<- int) {
		defer close(ch)

		for i := range data {
			ch <- i
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
	fmt.Println("###########################")
}

func lexicalConfinement() {
	fmt.Println("############## lexicalConfinement ##############")

	chanOwner := func() <-chan int {
		ch := make(chan int)

		go func() {
			defer close(ch)
			for i := 0; i < 4; i++ {
				ch <- i
			}
		}()

		return ch
	}

	chanConsumer := func(ch <-chan int) {
		for val := range ch {
			fmt.Printf("Received val %d \n", val)
		}
	}

	readStream := chanOwner()
	chanConsumer(readStream)

	fmt.Println("###########################")
}

func main() {
	// adhocConfinement()
	lexicalConfinement()
}
