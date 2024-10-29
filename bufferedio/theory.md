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



