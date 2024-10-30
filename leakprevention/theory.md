## Preventing Goroutine Leaks

- Goroutines are cheap and easy to create; it is one of the things that makes go a productive lanugage.

- The runtime handles multiplexing of goroutines onto any number of OS threads so we dont have to worry about abstraction.

- But they do cost resources and goroutines are not garbage collected by runtime.

- The gorutines has a few paths to termination
    - When it has completed its work
    - When it cannot continue its work due to an unrecoverable error
    - When it is told to stop working

- The first two paths are free, these parts are the algorithm but what about ``cancellation``?

- If a goroutine has begun it is most-likely cooperating with several other goroutines in somesort of organized fashion.

- We could even represent this interconnectedness as a graph whether a child goroutine must be execued or not can predicted by knowledge of other goroutines.

- The parent goroutiine (the main goroutine) must be able to tell the child goroutine to terminate with full contextual example.

## Goroutines and garbage collection

Goroutines in Go are, in fact, managed by the Go runtime, and they are garbage collected when they are no longer referenced or reachable. However, a goroutine must exit naturally for it to be eligible for garbage collection. The runtime keeps track of all running goroutines, and if a goroutine is "leaked" (for instance, if it's blocked indefinitely on a channel or an unending loop), it will not be garbage collected because it's still active.


-  **Completion**: If a goroutine completes its function or returns, it becomes eligible for garbage collection.
-  **Reachability**: If a goroutine is no longer reachable (e.g., no other goroutines can reference it), it becomes a candidate for garbage collection.
-  **Blocked/Leaked Goroutines**: If a goroutine is stuck or waiting indefinitely, it is not considered eligible for garbage collection, as the runtime treats it as still “alive.”

- To avoid goroutine leaks, make sure they have a clear exit strategy, especially when using channels or waiting for external input.

## Sample

```go
func sampleGoroutineLeak() {
	// doWork: this routine is not garbage collected and lives
	// whole lifetime of the program
	doWork := func(strs <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("do work exited")
			defer close(completed)
			for s := range strs {
				fmt.Println(s)
			}
		}()
		return completed
	}

	doWork(nil)
	fmt.Println("done!!")
}

func main() {
	sampleGoroutineLeak()
}
```
- The ``strs`` channel will not get any strings written into it and the goroutine containing ``doWork`` will be in memory for the lifetime of this process.

## Mitigation

- The way to mitigate this is to establish a signal between the parent goroutine and children that allows parent to ``signal cancellation`` to its children.

- The signal is ususally a ``read-only-channel`` named ``done``.

- The parent goroutine passes this channel to child goroutine and the parent goroutine closes this channel to close the child goroutine.

> If a goroutine is responsible for creating a goroutine it is also responsible for ensuring that it can stop the goroutine.

