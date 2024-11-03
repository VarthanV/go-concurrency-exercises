package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type Name struct {
	Name   string
	Locale string
}

func generateStreamFromFile(fileName string, locale string) <-chan Name {
	fileStream := make(chan Name)

	go func() {
		defer close(fileStream)
		file, err := os.Open(fileName)
		if err != nil {
			log.Println("error in opening file ", err)
			return
		}

		scanner := bufio.NewScanner(file)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			fileStream <- Name{
				Locale: locale,
				Name:   strings.ToUpper(scanner.Text()),
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	}()

	return fileStream
}

func fanIn(done <-chan interface{}, ch ...<-chan Name) <-chan Name {
	var (
		wg sync.WaitGroup
	)
	outStream := make(chan Name)

	wg.Add(len(ch))

	for _, c := range ch {
		go func(ch <-chan Name) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case val, ok := <-ch:
					if !ok {
						return
					}
					outStream <- val
				}
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(outStream)
	}()

	return outStream

}

func fanInDriver() {
	count := 0
	fmt.Println("################# fanInDriver ######################## ")
	done := make(chan interface{})
	defer close(done)

	spanishStream := generateStreamFromFile("spanish.txt", "es")
	britishStream := generateStreamFromFile("british.txt", "en")

	mergedStream := fanIn(done, spanishStream, britishStream)

	for val := range mergedStream {
		log.Printf("Name: %s , Locale: %s \n", val.Name, val.Locale)
		count += 1
	}
	log.Printf("Processed total %d names\n", count)

	fmt.Println("######################################### ")

}
