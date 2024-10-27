# WaitGroup

 - ``WaitGroup`` is a basically a way to wait for several goroutines to finish their work

 ## Typical WatiGroup flow

 - **Adding goroutines**: Before starting your goroutines, you tell the WaitGroup how many to expect. You do this with ``WaitGroup.Add(n)``, where n is the number of goroutines you’re planning to run.

 - **Goroutines running**: Each goroutine goes off and does its thing. When it’s done, it should let the WaitGroup know by calling ``WaitGroup.Done()`` to reduce the counter by one.

 - **Waiting for all goroutines**: In the main goroutine, the one not doing the heavy lifting, you call ``WaitGroup.Wait()``. This pauses the main goroutine until that counter in the WaitGroup reaches zero. In plain terms, it waits until all the other goroutines have finished and signaled they’re done.

```go
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        ...
    }()
}
```

## Why ``wg.Add(n)`` is error prone?

- The point is this, if the logic of the loop changes down the road, like if someone adds a ``continue`` statement that skips certain iterations, things can get messy


## Anatomy

```go
type WaitGroup struct {
	noCopy noCopy

	state atomic.Uint64
	sema  uint32
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
```

- In Go, it’s easy to copy a struct by just assigning it to another variable. But some structs, like ``WaitGroup``, really ``shouldn’t be copied``.

- Copying a WaitGroup can mess things up because the internal state that tracks the goroutines and their synchornization can get out of sync between the copie. 

## noCopy

- The ``noCopy`` struct is included in ``WaitGroup `` as a way to help prevent copying mistakes by throwing errors, but by serving as warning.

- The ``noCopy`` struct doesn’t actually affect how your program runs. Instead, it acts as a marker that tools like ``go vet`` can pick up on to detect when a struct has been copied in a way that it shouldn’t be.

```go
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
```

- Its structs is super simple
    - It has no fields, so it doesn’t take up any meaningful space in memory.
    - It has two methods, Lock and Unlock, which do nothing (no-op). These methods are there just to work with the`` -copylocks`` checker in the go vet tool.


- When you run go ``vet`` on your code, it checks to see if any structs with a ``noCopy`` field, like WaitGroup, have been copied in a way that could cause issues.



```go
func main() {
	var a sync.WaitGroup
	b := a

	fmt.Println(a, b)
}

// go vet:
// assignment copies lock value to b: sync.WaitGroup contains sync.noCopy
// call of fmt.Println copies lock value: sync.WaitGroup contains sync.noCopy
// call of fmt.Println copies lock value: sync.WaitGroup contains sync.noCopy
```

[Code](https://go.dev/play/p/8D42-xGo5jy)

## Internal State

- The state of a ``WaitGroup`` is stored in an ``atomic.Uint64`` variable

![Internal State](https://victoriametrics.com/blog/go-sync-waitgroup/sync-waitgroup-struct.webp)


- **Counter (high 32 bits)**: This part keeps track of the number of goroutines the WaitGroup is waiting for. When you call **wg.Add()** with a positive value, it bumps up this counter, and when you call **wg.Done()**, it **decreases** the counter by one.

- **Waiter (low 32 bits)**: This tracks the number of **goroutines currently waiting for that counter (the high 32 bits) to hit zero**. Every time you call wg.Wait(), it increases this **waiter** count. Once the counter reaches zero, it releases all the goroutines that were waiting.

- Then there’s the final field, ``sema uint32``, which is an internal semaphore managed by the Go runtime.

- When a goroutine calls`` wg.Wait()`` and the counter isn’t zero, it increases the waiter count and then blocks by calling ``runtime_Semacquire(&wg.sema)``. This function call puts the goroutine to sleep until it gets woken up by a corresponding ``runtime_Semrelease(&wg.sema) ``call.

## Alignment problem

- Evolution of WaitGroup struct

![alt text](https://victoriametrics.com/blog/go-sync-waitgroup/sync-waitgroup-versions.webp)

- When we talk about alignment, we’re referring to the need for data types to be stored at specific memory addresses to allow for efficient access.

- For example, on a 64-bit system, a ``64-bit`` value like ``uint64 ``should ideally be stored at a memory address that’s a multiple of 8 bytes. The reason is, the CPU can grab aligned data in one go, but if the data isn’t aligned, it might take multiple operations to access it.

![Alignment issues](https://victoriametrics.com/blog/go-sync-waitgroup/sync-waitgroup-alignment.webp)

- On ``32-bit architectures``, the compiler doesn’t guarantee that`` 64-bit`` values will be aligned on an ``8-byte boundary``. Instead, they might only be aligned on a ``4-byte boundary``.


