# Mutex

- Mutex or Mutual Exclusion is a Go way to make sure that only one goroutine is doing critical operation with a shared resource at a time.

- This resource can be a piece of code, an integer, a map, a struct, a channel, or pretty much anything.

- When a goroutine locks a mutex, it’s basically saying, ‘Hey, I’m going to use this shared resource for a bit,’ and every other goroutine has to wait until the mutex is unlocked. Once it’s done, it should unlock the mutex so other goroutines can get their turn.

```go
var counter = 0
var wg sync.WaitGroup

func incrementCounter() {
	counter++
	wg.Done()
}

func main() {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go incrementCounter()
	}

	wg.Wait()
	fmt.Println(counter)
}

```

- A ``race condition`` happens when multiple goroutines try to access and change shared data at the same time without proper synchronization. In this case, the increment operation (``counter++``) isn’t atomic.

- It’s made up of multiple steps, below is the Go assembly code for counter++ in ``ARM64 architecture``:

```bash
    MOVD	main.counter(SB), R0
    ADD	$1, R0, R0
    MOVD	R0, main.counter(SB)
```

- The counter++ is a ``read-modify-write`` operation and these steps above aren’t atomic, meaning they’re not executed as a ``single, uninterruptible action``.

- ![Race condition](https://victoriametrics.com/blog/go-sync-mutex/mutex-race-condition.webp)


- Solving using mutex

```go
var counter = 0
var wg sync.WaitGroup
var mutex sync.Mutex

func incrementCounter() {
	mutex.Lock()
	counter++
	mutex.Unlock()
	wg.Done()
}

func main() {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go incrementCounter()
	}

	wg.Wait()
	fmt.Println(counter)
}

```

## The Anatomy

