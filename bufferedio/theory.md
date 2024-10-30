# Buffered I/O

- Buffered I/O refers to the technique of temporarily storing the results of an I/O operation in the user-space before transmitting it to the kernel(in the case of write) or before providing it to your process in case of reads.

- By so buffering we can minimize the number of syscalls and can ``block -align`` I/O operations , which may improve the performace of your app. 

 - For example consider a process that writes one character at  a time to a file. This is inefficient , Each write operations corresponds to a ``write()`` syscall which means trip into the kernel , a memory copy (of a single byte) and return to user-space only to repeat the processs again.

 - Worst file systems and storage media work in terms of ``blocks``, operations are fastest when aligned to integer multiples of those blocks. Misaligned operations , particularly very small ones incur additional overhead.

 - User buffered I/O avoids this inefficienty by buffering the writes in a data buffer in a user space until a certain thershold is reached.

- Ideally the underlying filesystem's block size or an integer multiple therof.

- To use our previous example, we will simply copy each char into buffer and call ``write()`` only when the block size is reached.

- Similar process happens with reads, imagine a process that reads one line of a file into memory at time .

- The user buffered I/O library in C is called Standard I/O: ``fopen()`` to open a file, ``fwrite()`` to write, ``fread()`` to read, and so on.

``producer --> buffer --> io.Writer``

- Let’s visualise how buffering works with nine writes (one character each) when buffer has space for 4 characters:

```sh
producer         buffer           destination (io.Writer)
 
   a    ----->   a
   b    ----->   ab
   c    ----->   abc
   d    ----->   abcd
   e    ----->   e      ------>   abcd
   f    ----->   ef               abcd
   g    ----->   efg              abcd
   h    ----->   efgh             abcd
   i    ----->   i      ------>   abcdefgh

```

- ``bufio.Writer`` uses ``[]byte`` buffer under the hood 

```go
type Writer int
func (*Writer) Write(p []byte) (n int, err error) {
    fmt.Println(len(p))
    return len(p), nil
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
}
```

- Unbuffered I/O simply means that each write operation goes straight to destination. We’ve 4 write operations and each one maps to ``Write`` call where passed slice of bytes has length 1.

- With buffered I/O we’ve internal buffer (3 bytes long) which collects data and ``flushes`` buffer when full. First three writes end up inside the buffer.

- 4th write detects buffer with no free space so it sends accumulate data out.

- It gives space to hold ``d`` character, theres something more - Flush call. its needed at very end to flush any outstanding data.

- ``bufio.Writer`` sends data only when buffer is either full or when explicitly requested with ``Flush`` method.

-  By default bufio.Writer uses ``4096`` bytes long buffer. It can be set with ``NewWriterSize`` function.

## Implementation

- It’s rather straightforward

```go
type Writer struct {
    err error
    buf []byte
    n   int
    wr  io.Writer
}

```

- Field ``buf`` accumlates data , Consumer(wr) gets data when buffer is full or ``Flush`` is called.

- First encountered I/O error is held by ``err`` after encountering an error , writer is no-op

```go
type Writer int
func (*Writer) Write(p []byte) (n int, err error) {
    fmt.Printf("Write: %q\n", p)
    return 0, errors.New("boom!")
}
func main() {
    w := new(Writer)
    bw := bufio.NewWriterSize(w, 3)
    bw.Write([]byte{'a'})
    bw.Write([]byte{'b'})
    bw.Write([]byte{'c'})
    bw.Write([]byte{'d'})
    err := bw.Flush()
    fmt.Println(err)
}
Write: "abc"
boom!
```


- Here we see that ``Flush`` didn’t call 2nd write on our consumer. Buffered writer simply doesn’t try to do more writes after first error.


- Field `n` is the ``current writing position`` inside the buffer. Buffered method returns n’s value

```go
type Writer int
func (*Writer) Write(p []byte) (n int, err error) {
    return len(p), nil
}
func main() {
    w := new(Writer)
    bw := bufio.NewWriterSize(w, 3)
    fmt.Println(bw.Buffered())
    bw.Write([]byte{'a'})
    fmt.Println(bw.Buffered())
    bw.Write([]byte{'b'})
    fmt.Println(bw.Buffered())
    bw.Write([]byte{'c'})
    fmt.Println(bw.Buffered())
    bw.Write([]byte{'d'})
    fmt.Println(bw.Buffered())
}
0
1
2
3
```

- It starts with 0 and is incremented by the number of bytes added to buffer. It’s also reset after flush to underlying writer while calling

```go 
    bw.Write([]byte{'d'})
 ```

## Large Writes

```go
    type Writer int
    func (*Writer) Write(p []byte) (n int, err error) {
        fmt.Printf("%q\n", p)
        return len(p), nil
    }
    func main() {
        w := new(Writer)
        bw := bufio.NewWriterSize(w, 3)
        bw.Write([]byte("abcd"))
    }
```

- prints "**abcd"** because ``bufio.Writer`` detects if ``Write`` is called with amount of data ``too much for internal buffer`` (3 bytes in this case) . It then calls ``Write`` method directly on writer (destination object). It’s completely fine since amount of data is large enough to skip proxying through temporary buffer.

## Reset

- Buffer which is the core part of bufio.Writer can be re-used for different destination writer with ``Reset`` method. It saves memory allocation and extra work for garbage collector 

```go
    type Writer1 int
    func (*Writer1) Write(p []byte) (n int, err error) {
        fmt.Printf("writer#1: %q\n", p)
        return len(p), nil
    }
    type Writer2 int
    func (*Writer2) Write(p []byte) (n int, err error) {
        fmt.Printf("writer#2: %q\n", p)
        return len(p), nil
    }
    func main() {
        w1 := new(Writer1)
        bw := bufio.NewWriterSize(w1, 2)
        bw.Write([]byte("ab"))
        bw.Write([]byte("cd"))
        w2 := new(Writer2)
        bw.Reset(w2)
        bw.Write([]byte("ef"))
        bw.Flush()
    }
    writer#1: "ab"
    writer#2: "ef"
```

- There is one bug in this program. Before calling Reset we should flush the buffer with Flush. Currently, written data cd is lost since Reset simply discards any outstanding information.

## Buffer free space

- To check how much space left inside the buffer we can use ``Available`` method.

```go
    w := new(Writer)
    bw := bufio.NewWriterSize(w, 2)
    fmt.Println(bw.Available())
    bw.Write([]byte{'a'})
    fmt.Println(bw.Available())
    bw.Write([]byte{'b'})
    fmt.Println(bw.Available())
    bw.Write([]byte{'c'})
    fmt.Println(bw.Available())
    2
    1
    0
    1
```

## Write{Byte,Rune,String} Methods

- At our disposal we’ve 3 utility functions to write data of common types

```go
    w := new(Writer)
    bw := bufio.NewWriterSize(w, 10)
    fmt.Println(bw.Buffered())
    bw.WriteByte('a')
    fmt.Println(bw.Buffered())
    bw.WriteRune('ł') // 'ł' occupies 2 bytes
    fmt.Println(bw.Buffered())
    bw.WriteString("aa")
    fmt.Println(bw.Buffered())
    0
    1
    3
    5
```

## ReadFrom

- Package io defines ``io.ReaderFrom`` interface. It’s usually implemented by writer to do the dirty work of reading all the data from specified reader (until EOF):

```go
type ReaderFrom interface {
        ReadFrom(r Reader) (n int64, err error)
}
```

- ``bufio.Writer`` implements this interface, allowing to call ``ReadFrom`` method which ``digests all data`` from io.Reader

```go
type Writer int
func (*Writer) Write(p []byte) (n int, err error) {
    fmt.Printf("%q\n", p)
    return len(p), nil
}
func main() {
    s := strings.NewReader("onetwothree")
    w := new(Writer)
    bw := bufio.NewWriterSize(w, 3)
    bw.ReadFrom(s)
    err := bw.Flush()
    if err != nil {
        panic(err)
    }
}
"one"
"two"
"thr"
"ee"
```
- It’s important to call ``Flush`` even while using ReadFrom.

## bufio.Reader

- It allows to read in bigger batches from the underlying ``io.Reader``.

-  This leads to less read operations which can improve performance if e.f. underlying media works better when data is read in blocks of certain size

```
io.Reader --> buffer --> consumer

```