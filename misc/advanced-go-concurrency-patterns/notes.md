# Advanced Go Concurrency Patterns - Sameer Ajmani 

[Talk](https://www.youtube.com/watch?v=QDDwwePbDtw&t=1s)


## Goroutines and channels 

- Goroutines are independenly executing functions in the same address space.

```go
go f()
go g(1,2)
```
- Channels are type values that allow goroutines to synchornize and exchange information

```go
c := make(chan int)

go func(){c <- 3}()
n := <-c
```

- For more details check on  [basics](https://github.com/VarthanV/go-concurrency-exercises/tree/main/misc/google-io-2012)

## Basic Ping Pong

```go
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
```

## Deadlock detection

```go
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

	// table <- new(Ball) // game on ,toss the ball
	time.Sleep(1 * time.Second)
	<-table //game over grab the ball
	fmt.Println("#######################")
}

func main() {
	pingPong()
}

```

```sh
#########  pingPong ###########
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan receive]:
main.pingPong()
        /Users/varthanv/Desktop/go-concurrency-exercises/misc/advanced-go-concurrency-patterns/main.go:28 +0x111
main.main()
        /Users/varthanv/Desktop/go-concurrency-exercises/misc/advanced-go-concurrency-patterns/main.go:33 +0xf

goroutine 18 [chan receive]:
main.pingPong.func1({0x37691c0, 0x4}, 0xc00008e0c0)
        /Users/varthanv/Desktop/go-concurrency-exercises/misc/advanced-go-concurrency-patterns/main.go:15 +0x37
created by main.pingPong in goroutine 1
        /Users/varthanv/Desktop/go-concurrency-exercises/misc/advanced-go-concurrency-patterns/main.go:23 +0xaa

goroutine 19 [chan receive]:
main.pingPong.func1({0x37691bc, 0x4}, 0xc00008e0c0)
        /Users/varthanv/Desktop/go-concurrency-exercises/misc/advanced-go-concurrency-patterns/main.go:15 +0x37
created by main.pingPong in goroutine 1
        /Users/varthanv/Desktop/go-concurrency-exercises/misc/advanced-go-concurrency-patterns/main.go:24 +0xf6
exit status 2
```

## Panic dumps the stack

```go
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
    panic("need the stacks")
}

func main() {
	pingPong()
}

```

```sh
#########  pingPong ###########
pong 1
ping 2
pong 3
ping 4
pong 5
ping 6
pong 7
ping 8
pong 9
ping 10
pong 11
#######################
panic: need the stacks

goroutine 1 [running]:
main.pingPong()
        /Users/varthanv/Desktop/go-concurrency-exercises/misc/advanced-go-concurrency-patterns/main.go:30 +0x165
main.main()
        /Users/varthanv/Desktop/go-concurrency-exercises/misc/advanced-go-concurrency-patterns/main.go:34 +0xf
exit status 2
```

- We could see one goroutine getting leaked in the above stack trace

## Its easy to go but how to stop?

- Long-lived programs need to cleanup

- Need to write programs that handle communication, periodic events and cancellation

- The core is Go's ``select`` statement; like ``switch`` but the decision is made based on the ability to communicate

```go
select {
    case xc <- x :
    // sent x on xc
    case y := <-yc:
    // received y from yc
}
```

## Race detection using cli

```sh
go run -race <file.go>
```

## Structure: for-select loop 

- ``loop`` runs its own goroutine

- ``select`` lets ``loop`` avoid blocking indefinitely in any one state

```go
func (s*sub) loop(){
    // declare mutabke state
    for {
        select {
            case <-c1:
                // read-/write state
            case c2<- x:
                // read/write state
            case y := <-c3
                // read write state
        }
    }
}
```

- The cases interact via local state in the ``loop``

## Select and nil channels:

- Sends and receives on a nil channels block

- Select never selects a blocking case


# Slides

[Slides](https://go.dev/talks/2013/advconc.slide#1)

