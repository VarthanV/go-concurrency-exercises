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

- At its core a mutex in Go has two fields; ``state`` and ``sema``.

```go
package sync

type Mutex struct {
    state int32

}
```

- The ``state`` field is a 32-bit integer that shows the current state of the mutex. It is divided into multiple bits that encode various piece of information about the mutex.

![Mutex structure](https://victoriametrics.com/blog/go-sync-mutex/mutex-structure.webp)

- Below is the explanation of the state

    - **Locked(bit 0)**: Whether the mutex is currently locked, If it set to 1 the mutext is locked and no other goroutine can grab it.

    - **Woken (bit 1)**: Set to 1 if any goroutine has been woken up and is trying to acquire mutex. Other goroutines shouldn't be woken up unnecessarily.

    - **Starvation(bit 2)**: The bit shows whether mutex is in starvation mode(set to 1)

    - **Waiter (bit 3-31)**: The rest of the bits keep track of how many goroutines are waiting for the mutex.

- The other field ``sema``, is a ``uint32`` that acts as semaphore to manage and signal waiting goroutines. 

- When the mutex is unlocked , one of the waiting goroutines is woken up to acquire the lock

## Mutex Lock flow

- In the ``mutex.Lock`` fn , there are two paths; the ``fast path`` for the usual case and ``slow path`` for handling the unusual case.

```go
func (m *Mutex) Lock() {
	// Fast path: grab unlocked mutex.
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// Slow path (outlined so that the fast path can be inlined)
	m.lockSlow()
}
```

-  The fast path is designed to be really quick and is expected to handle ``most lock acquisitions`` where the mutex isn’t already in use.

- This path is also inlined, meaning it’s embedded directly into the calling function.

```bash
$ go build -gcflags="-m"

./main.go:13:12: inlining call to sync.(*Mutex).Lock
./main.go:15:14: inlining call to sync.(*Mutex).Unlock
```
 

- When the ``CAS (Compare And Swap)`` operation in the fast path fails, it means the state field wasn’t 0, so the mutex is currently locked.

- The real concern here is the slow path m.lockSlow, which does most of the heavy lifting.

- In the slow path ,the goroutine keeps actively spinning to ``try to acquire the lock``. It doesn't go straight to the waiting queue.

- Spinning means the goroutine enters a tight loop , Repeatedly checking the state of the mutex without giving up the CPU

- In this case, it is not a simple for loop but low-level assembly instructions to perform the ``spin-wait``.

```bash
TEXT runtime·procyield(SB),NOSPLIT,$0-0
	MOVWU	cycles+0(FP), R0
again:
	YIELD
	SUBW	$1, R0
	CBNZ	R0, again
	RET

```

- The assembly code runs a tight loop for 30cycles (``runtime.procyield(30)``) repeatedly yielding the CPU and decrementing the spin counter.

- The idea behind spinning is to wait for a short while in the hopes that mutex will free up soon. Letting the goroutine grab the mutext without the overhead of ``sleeo-wake-cycle``.

- In a ``single core`` machine spinning isnt enabled because it would just waste CPU time.

> “But what if another goroutine is already waiting for the mutex? It doesn’t seem fair if this goroutine takes the lock first.”


- That's why our mutex has two modes
     - Normal mode
     - Starvation mode (Spinning disabled)

## Normal mode
-  In normal mode , goroutines waiting for the mutex are organized in a first-in,first-out(FIFO) queue.

- When a goroutines wakes up to try and grab the mutex it doesn't get control immediately. Instead, it has to compete with any new goroutines that also want mutex at that time.

- The competition is titled in favor of new gorutines because theyr'e already running on CPU and can quickly try t grab the mutex , while the queued goroutine is still waking up.

![Normal mode](https://victoriametrics.com/blog/go-sync-mutex/mutex-normal-mode.webp)

> “What if that goroutine is unlucky and always wakes up when a new goroutine arrives?”

## Stravation mode 

- If that happens , it never acquires the lock, That's why we need to switch the mutex into ``starvation mode``.

- Stravation mode kicks in if a goroutine fails to acquire the lock for more than 1 ms . Its designed to make sure that waiting goroutines eventually get a ``fair chance`` at the mutex.

- In this mode, when a goroutine releases the mutex, it directly passes control to the goroutine at the front of the queue. This means no competition, no race, from new goroutines. They don’t even try to acquire it and just join the end of the waiting queue.

![Mutex starvation](https://victoriametrics.com/blog/go-sync-mutex/mutex-starvation-mode.webp)

- In the image above, the mutex continues giving the access to G1, G2, and so on. Each goroutine that has been waiting gets control and checks two conditions:

- If it is the last goroutine in the waiting queue.
- If it had to wait for less than one millisecond.
- If either of these conditions is true, the mutex switches back to normal mode.

## Mutex Unlock flow

- In unlock flow we have two paths; the fast path which is inlined and the slow path which handles unusual case.

- The fast path drops the locked bit in the state of mutex.

- If dropping this bit makes the state zero, it means no other flags are set (like waiting goroutines), and our mutex is now completely free.

- That’s where the slow path comes in and it needs to know if our mutex is in normal mode or starvation mode. Here’s a look at the slow path implementation:

```go
func (m *Mutex) unlockSlow(new int32) {
	// 1. Attempting to unlock an already unlocked mutex
	// will cause a fatal error.
	if (new+mutexLocked)&mutexLocked == 0 {
		fatal("sync: unlock of unlocked mutex")
	}
	if new&mutexStarving == 0 {
		old := new
		for {
			// 2. If there are no waiters, or if the mutex is already locked,
			// or woken, or in starvation mode, return.
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
				return
			}
			// Grab the right to wake someone.
			new = (old - 1<<mutexWaiterShift) | mutexWoken
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				runtime_Semrelease(&m.sema, false, 1)
				return
			}
			old = m.state
		}
	} else {
		// 3. If the mutex is in starvation mode, hand off the ownership
		// to the first waiting goroutine in the queue.
		runtime_Semrelease(&m.sema, true, 1)
	}
}

```


- In normal mode, if there are waiters and no other goroutine has been woken or acquired the lock, the mutex tries to decrement the waiter count and turn on the ``mutexWoken`` flag atomically.

- In starvation mode, it atomically increments the semaphore (``mutex.sem``) and hands off mutex ownership directly to the first waiting goroutine in the queue. The second argument of ``runtime_Semrelease`` determines if the handoff is ``true``.


## Inspired from

[go-sync-mutex](https://victoriametrics.com/blog/go-sync-mutex/index.html)