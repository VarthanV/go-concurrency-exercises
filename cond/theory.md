## Cond

- In golang ``sync.Cond`` is a synchornization primitive

- When a goroutine needs to wait for something specific to happen like some shared data changing it can ``block`` , meaning it just pauses its work until it gets a go-ahead to continue.

- The most basic way to do this is with a loop may evenadding ``time.Sleep`` to prevent the CPU from busy-waiting

```go
// wait until condition is true
for !condition {
}

// or
for !condition {
    time.Sleep(100 * time.Millisecond)
}
```

- This is not really efficient as the loop is still running in the background burning through CPU cycles even nothing is changed.

- When one goroutine is waiting for something to happen (waiting for a certain condition to become true), it can call ``Wait()``.

-  Another goroutine, once it knows that the condition might be met, can call ``Signal()`` or ``Broadcast()`` to wake up the waiting goroutine(s) and let them know it’s time to move on.

- Below is the basic interface that ``sync.Cond`` provides

```go
// Suspends the calling goroutine until the condition is met
func (c *Cond) Wait() {}

// Wakes up one waiting goroutine, if there is one
func (c *Cond) Signal() {}

// Wakes up all waiting goroutines
func (c *Cond) Broadcast() {}
```

![Overview](https://victoriametrics.com/blog/go-sync-cond/go-sync-cond-overview.webp)

- Alright, let’s check out a quick pseudo-example. This time, we’ve got a Pokémon theme going on, imagine we’re waiting for a specific Pokémon, and we want to notify other goroutines when it shows up.



```go
var pokemonList = []string{"Pikachu", "Charmander", "Squirtle", "Bulbasaur", "Jigglypuff"}
var cond = sync.NewCond(&sync.Mutex{})
var pokemon = ""

func main() {
	// Consumer
	go func() {
		cond.L.Lock()
		defer cond.L.Unlock()

		// waits until Pikachu appears
		for pokemon != "Pikachu" {
			cond.Wait()
		}
		println("Caught" + pokemon)
		pokemon = ""
	}()

    // Producer
	go func() {
		// Every 1ms, a random Pokémon appears
		for i := 0; i < 100; i++ {
			time.Sleep(time.Millisecond)

			cond.L.Lock()
			pokemon = pokemonList[rand.Intn(len(pokemonList))]
			cond.L.Unlock()

			cond.Signal()
		}
	}()

	time.Sleep(100 * time.Millisecond) // lazy wait
}

// Output:
// Caught Pikachu
```

- The problem is, there’s a gap between the producer sending the signal and the consumer actually waking up. In the meantime, the Pokémon could change, because the consumer goroutine might wake up later than 1ms (rarely) or other goroutine modifies the shared pokemon. So sync.Cond is basically saying: ‘Hey, something changed! Wake up and check it out, but if you’re too late, it might change again.’

- In fact, ``channels`` are generally preferred over ``sync.Cond`` in Go because they’re simpler, more idiomatic, and familiar to most developers.


- There is Github issue to remove the Cond type also

[GH issue](https://github.com/golang/go/issues/21165)

## Best usage

-  Aconsistent pattern in consumer: we always ``lock`` the mutex before waiting (``.Wait()``) on the condition, and we ``unlock`` it after the condition is met.

- Plus we wrap the condition inside a loop

```go
// Consumer
go func() {
	cond.L.Lock()
	defer cond.L.Unlock()

	// waits until Pikachu appears
	for pokemon != "Pikachu" {
		cond.Wait()
	}
	println("Caught" + pokemon)
}()
```

## Cond.Wait

- When we call ``Wait`` on ``sync.Cond`` we are telling the current goroutine to hang tight until the condition is met

- Heres what happens behind the  scenes
    - A goroutine gets added to the list of other goroutines that are also waiting on this same condition. All these goroutines are blocked , meaning they cant be woken up either by a ``Signal`` or ``Broadcast`` call.
    - The key part here is the mutet must be locked before calling ``Wait`` because ``Wait`` doess something important it automatically release the lock call before putting the goroutine to sleeo (defer ``Unlock``) this allows other goroutine to grab lock and do their work while the original goroutine is waitng.

    - When the waiting goroutine gets woken up (by ``Signal() ``or ``Broadcast()``), it doesn’t immediately resume work. First, it has to re-acquire the lock (``Lock()``).

![Wait call](https://victoriametrics.com/blog/go-sync-cond/go-sync-cond-wait.webp)

```go
func (c *Cond) Wait() {
	// Check if Cond has been copied
	c.checker.check()

	// Get the ticket number
	t := runtime_notifyListAdd(&c.notify)

	// Unlock the mutex
	c.L.Unlock()

	// Suspend the goroutine until being woken up
	runtime_notifyListWait(&c.notify, t)

	// Re-lock the mutex
	c.L.Lock()
}
```

## Further reading


[Sync-cond](https://victoriametrics.com/blog/go-sync-cond/)