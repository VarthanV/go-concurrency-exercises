# Confinement pattern

- When working with concurrent code there are few different options for safe-operation
    - Synchornization primitives for sharing memory (eg:``sync.Mutex``)
    
    - Sychonrnization via communication (eg,``channels``)

- There are couple of options that are implicitly safe within multiple concurrent processes
    - Immutable data.
    
    - Data protected by confinement.


## Immutable data

- In some case immutable data is ideal because it implicitly is ``concurrent-safe``.

- Each concurrent process may operate on the same data, but it may not modify it.

- If it wants to create the new data , it must create a new copy of the data with desired modifications.

- This allows not only a lighter cognitive load on developer but can also lead to faster programs if it leads to smaller critical sections or eliminates them altogether.

- In Go , we can achieve this by writing code that utilizes copies of values instead of pointers to values in memory.

## Confinement

- Confinement can also allow for lighter cognitive load on developer and smaller critical sections.

- The techinques to confine concurrent values are bit more involved than simply passing copies of value.

- It is a simple yet powerful idea of ensuring information is only ever available from ``one concunrrent`` process.

- When this is acheieved a concurrent program is ``implicitly`` safe and ``no synchornization`` needed.

- There are two kinds of confinement
    - Adhoc
    - Lexical

## Adhoc confinement

- It is confinement which is achieved through a convention whether it is set by a language community or group youb work within.

- It is difficult to stick just by word of mouth when we are working across a larger organisation , we need to do static check on everytime we commit.

```go
func adhocConfinement() {
	fmt.Println("############## adhocConfinement ##############")

	data := make(chan int, 4)

	loopData := func(ch chan<- int) {
		defer close(ch)

		for i := range data {
			ch <- i
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
	fmt.Println("###########################")
}
```

## Lexical Confinement

- It involves using lexical scope to expose only the correct data and conncurrency primitives for multiple concurrent process to use.

-  It makes impossible to do wrong thing.

```go
func lexicalConfinement() {
	fmt.Println("############## lexicalConfinement ##############")

	chanOwner := func() <-chan int {
		ch := make(chan int)

		go func() {
			defer close(ch)
			for i := 0; i < 4; i++ {
				ch <- i
			}
		}()

		return ch
	}

	chanConsumer := func(ch <-chan int) {
		for val := range ch {
			fmt.Printf("Received val %d \n", val)
		}
	}

	readStream := chanOwner()
	chanConsumer(readStream)

	fmt.Println("###########################")
}

```

- Synchornization comes with a cost , if we can avoid it we wont have any critical sections and therefore we dont have to pay the cost for synchorizing.

- Concurrenrt code that utilizies lexical confinement is usually simpler to understand than the concurrent code without lexically confined varibales.

- This is because within the context of your ``lexical scope`` we can write synchornous code.

