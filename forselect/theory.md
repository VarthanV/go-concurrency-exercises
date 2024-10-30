## for-select loop

- Basic syntax

```go
    for { // Eithr loop infinitely or range over something
        select{
            // do some work with channels
        }
    }
```

- There are couple of scenarios where this pattern pops up


- In Go's `select` statement, if none of the specified channels are ready, and there’s no `default` case, the `select` statement will block and wait for one of the cases to become ready. However, if there is a `default` case, `select` will not block. Instead, it will immediately execute the `default` case if no other channel is ready, making it non-blocking.

When it comes to preemption:

- The `default` case in a `select` statement is **not preemptible** in the sense that it will execute immediately if no other channel operations are ready, without yielding to the scheduler to wait on other goroutines or channels.
  
- Without a `default` case, `select` will block and wait, allowing other goroutines to potentially proceed (and thus can be considered more cooperative with Go's scheduler, enabling preemption more naturally in the waiting state).

- In short, a `select` with a `default` case makes the code non-blocking and therefore can bypass waiting, while a blocking `select` (without `default`) is more preemption-friendly because it lets the goroutine yield.

## Usecases

- Sending iteration variables out on a channel

```go
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
```

- Loop infinitely waiting to be stopped

