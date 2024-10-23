package main

import (
	"fmt"
	"sync"
)

func main() {

	counterSynchornizedExample := func() {
		fmt.Println("########## counterSynchornizedExample ############# ")
		var wg sync.WaitGroup

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				fmt.Println(i)
			}(i)
		}

		wg.Wait()
		fmt.Println("####################### ")
	}

	// Closure functions of goroutines operates on the memory ref
	closuredExample := func() {
		fmt.Println("########## closuredExample ############# ")
		var wg sync.WaitGroup
		salutation := "hello!"

		wg.Add(1)
		go func() {
			defer wg.Done()
			salutation = "hola!!"
		}()
		wg.Wait()
		fmt.Println(salutation) // hola!!
		fmt.Println("####################### ")

	}

	/*
	 Closured function with forloop gotcha,
	 The function closes over the iteration variable salutation

	*/
	closuredExampleWithForLoopGotcha := func() {
		var wg sync.WaitGroup
		fmt.Println("########## closuredExampleWithForLoopGotcha ############# ")
		for _, salutation := range []string{"Hello", "Hola", "Vanakkam"} {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Println(salutation) // non deterministic output
			}()
		}

		wg.Wait()
		fmt.Println("####################### ")
	}

	closuredExampleForLoopGotchaFixed := func() {
		fmt.Println("########## closuredExampleForLoopGotchaFixed ############# ")
		var wg sync.WaitGroup

		for _, salutation := range []string{"Hello", "Hola", "Vanakkam"} {
			wg.Add(1)
			go func(s string) {
				defer wg.Done()
				fmt.Println(s)
			}(salutation)
		}

		wg.Wait()
		fmt.Println("####################### ")
	}

	closuredExampleWithForLoopGotcha()
	closuredExampleForLoopGotchaFixed()
	counterSynchornizedExample()
	closuredExample()
}
