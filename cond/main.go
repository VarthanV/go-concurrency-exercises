package main

import (
	"fmt"
	"sync"
	"time"
)

func buttonClickEventHandler() {
	type Button struct {
		Clicked *sync.Cond
	}

	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(cond *sync.Cond, fn func()) {
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			wg.Done()
			cond.L.Lock()
			defer cond.L.Unlock()
			cond.Wait()
			fn()
		}()

		wg.Wait()
	}

	var clickRegistered sync.WaitGroup

	clickRegistered.Add(2)

	subscribe(button.Clicked, func() {
		fmt.Println("click event 1")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("button click done")
		clickRegistered.Done()

	})

	button.Clicked.Broadcast()
	clickRegistered.Wait()
}

func main() {

	basicExample := func() {
		c := sync.NewCond(&sync.Mutex{})
		queue := make([]interface{}, 0, 10)

		removeFromQueue := func(delay time.Duration) {
			time.Sleep(delay)
			c.L.Lock()

			queue = queue[1:]
			fmt.Println("Removed from queue")
			c.L.Unlock()
			c.Signal()
		}

		for i := 0; i < 10; i++ {
			c.L.Lock()
			for len(queue) == 2 {
				c.Wait()
			}

			fmt.Println("adding to queue")
			queue = append(queue, struct{}{})
			go removeFromQueue(1 * time.Second)
			c.L.Unlock()
		}
	}

	basicExample()
	buttonClickEventHandler()

}
