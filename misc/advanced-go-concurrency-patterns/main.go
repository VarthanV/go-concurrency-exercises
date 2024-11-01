package main

import (
	"fmt"
	"time"
)

func pingPong() {
	fmt.Println("#########  pingPong ###########")
	type Ball struct{ hits int }
	table := make(chan *Ball)

	play := func(name string, table chan *Ball) {
		for {
			ball := <-table
			ball.hits++
			fmt.Println(name, ball.hits)
			time.Sleep(100 * time.Millisecond)
			table <- ball
		}
	}

	go play("ping", table)
	go play("pong", table)

	table <- new(Ball) // game on ,toss the ball
	time.Sleep(1 * time.Second)
	<-table //game over grab the ball
	fmt.Println("#######################")
}

func main() {
	pingPong()
}
