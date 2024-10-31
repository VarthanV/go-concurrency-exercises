package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Writer int

func (*Writer) Write(p []byte) (n int, err error) {
	fmt.Println(len(p))
	return len(p), nil
}

type fileWriter struct {
	f *os.File
}

func (fr *fileWriter) Write(p []byte) (n int, err error) {
	fmt.Println("Writing.. ", string(p))
	n, err = fr.f.Write(p)
	return n, err
}

func bufferedFileWriter() {

	done := make(chan interface{})
	defer close(done)

	f, err := os.Create("test.txt")
	if err != nil {
		log.Fatal("unable to open file ", err)
	}
	defer f.Close()

	w := &fileWriter{f: f}

	bw := bufio.NewWriterSize(w, 10)
	defer func() {
		err = bw.Flush()
		if err != nil {
			log.Fatal("unable to flush buffer ", err)
		}
	}()

	producer := func(done <-chan interface{}) <-chan string {
		ch := make(chan string)

		go func() {
			defer close(ch)
			for _, val := range []string{"foo", "bar", "bas", "fjffjfjfjfjfj"} {
				select {
				case ch <- val:

				case <-done:
					return
				}

			}
		}()
		return ch
	}

	for val := range producer(done) {
		_, err := bw.Write([]byte(val))
		if err != nil {
			log.Fatal("error in writing value ", err)
		}
	}

}

func main() {
	fmt.Println("Unbuffered I/O")
	w := new(Writer)
	w.Write([]byte{'a'})
	w.Write([]byte{'b'})
	w.Write([]byte{'c'})
	w.Write([]byte{'d'})
	fmt.Println("Buffered I/O")
	bw := bufio.NewWriterSize(w, 3)
	bw.Write([]byte{'a'})
	bw.Write([]byte{'b'})
	bw.Write([]byte{'c'})
	bw.Write([]byte{'d'})
	err := bw.Flush()
	if err != nil {
		panic(err)
	}
	bufferedFileWriter()
}
