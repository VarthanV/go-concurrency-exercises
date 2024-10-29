# Select 

- It is a glue that binds channels together

- It is how we are able to compose channels in a program to form larger abstractions.

- select statement binds together channels locally within a single function or type and also globally at intersection of two or more components in a system.

- select statement can help safely bring channels together with concepts like ``cancellation,timeouts,waiting and default values``.

```go
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
```

- Unlike ``switch`` blocks case statement in a select block are ``not tested sequentially``

- Execution will not fall through if none of the criteria are met.

- All channel reads and writes are considered simultaneously to see if any of them are ready (populated or closed channels in case of reads and channels that are not at capacity in case of writes).

- If none of the channels are ready the entire select statement blocks.

- Empty switch statements block forever

```go
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
```

## GOMAXPROCS Lever

- In the ``runtime`` package there is a function called ``GOMAXPROCS``.

- This fn controls the number of ``OS threads`` that will so host so-called ``worker queues``

```go
runtime.GOMAXPROCS(runtime.NumCPU())
```


