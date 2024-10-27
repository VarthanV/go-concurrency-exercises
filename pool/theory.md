## Pool

## Background

- It is commonly used in the standard library. For instance, in the ``encoding/json package``

```go
package json

var encodeStatePool sync.Pool

// An encodeState encodes JSON into a bytes.Buffer.
type encodeState struct {
	bytes.Buffer // accumulated output

	ptrLevel uint
	ptrSeen  map[any]struct{}
}
```


- In this case, ``sync.Pool`` is being used to reuse ``*encodeState`` objects, which handle the process of ``encoding JSON into a bytes.Buffer.``


## Basics

- ``sync.Pool`` in Go is a place where we can keep temporary objects for later reuse.

- We cannot control how many objects stay in the pool and anything we put in there can be removed at any time without any warning.

- The pool is built to be ``thread-safe`` , so multiple goroutines can tap into it simulataneously.

## Why bother reusing objects? 

- When we have got a lot of goroutines running at once , they often need similar objects, Imaging running ``go f()`` multiple times concurrently.

- If each goroutine creates is own objects memory usage can quickly increase and this puts a strain on the ``garbage-collector`` because it has to clean up all those objects once they are no longer needed.

- The stiuation creates a cycle where ``high concurrency leads to high memory usage`` which then slows down the garbage collector.

- ``sync.Pool`` is designed to help break this cycle.

```go
type Object struct {
	Data []byte
}

var pool sync.Pool = sync.Pool{
	New: func() any {
		return &Object{
			Data: make([]byte, 0, 1024),
		}
	},
}
```

-  A ``New()`` function that returns a new object when the pool is empty can be used to build a pool. If you omit this optional function, the pool just returns ``nil`` if it is empty.

- In the snippet above, the goal is to reuse the ``Object`` struct instance, specifically the slice inside it.


- Reusing the slice helps reduce unnecessary growth.

-  For instance if the slice grows to 8192 bytes during use , you can reset its length to zero before putting back into the pool . The undelrying arr still has capacity of 8192, So the next we time we need it , those 8192 bytes are ready to be reused.

```go
func (o *Object) Reset() {
	o.Data = o.Data[:0]
}

func main() {
	testObject := pool.Get().(*Object)

	// do something with testObject

	testObject.Reset()
	pool.Put(testObject)
}
```

- If we don't like using type assertions ``pool.Get().(*Object)``, there are couple of ways to avoid it.

- Use dedicated fn to get the object from the pool.

```go
func getObjectFromPool() (*Object,error) {
	obj ,ok := pool.Get().(*Object)
    if !ok {
        return nil,errors.New("error in getting object from pool")
    }
	return obj, nil
}
```

- Create generic version of `sync.Pool`

```go
type Pool[T any] struct {
	sync.Pool
}

func (p *Pool[T]) Get() T {
	return p.Pool.Get().(T)
}

func (p *Pool[T]) Put(x T) {
	p.Pool.Put(x)
}

func NewPool[T any](newF func() T) *Pool[T] {
	return &Pool[T]{
		Pool: sync.Pool{
			New: func() interface{} {
				return newF()
			},
		},
	}
}
```

- Just note that, it adds a tiny bit of overhead due to the ``extra layer of indirection``. In most cases, this overhead is minimal, but if you’re in a highly CPU-sensitive environment, it’s a good idea to run benchmarks to see if it’s worth it.

## Allocation trap

- what we store in the pool is typically not the object itself but a pointer to the object.

```go
var pool = sync.Pool{
	New: func() any {
		return []byte{}
	},
}

func main() {
	bytes := pool.Get().([]byte)

	// do something with bytes
	_ = bytes

	pool.Put(bytes)
}
```

- We’re using a pool of ``[]byte``. Generally (though not always), when you pass a value to an ``interface``, it may cause the value to be placed on the ``heap``. This happens here too, not just with slices but with anything you pass to ``pool.Put()`` that isn’t a pointer.

```sh
// escape analysis
$ go build -gcflags=-m

bytes escapes to heap
```

- However, if we pass a pointer to ``pool.Put()``, there is no extra allocation

```go
var pool = sync.Pool{
	New: func() any {
		return new([]byte)
	},
}

func main() {
	bytes := pool.Get().(*[]byte)

	// do something with bytes
	_ = bytes

	pool.Put(bytes)
}
```
- Sample example program [Pool](https://github.com/golang/go/blob/2580d0e08d5e9f979b943758d3c49877fb2324cb/src/sync/example_pool_test.go#L15)

## Internals

- PMG stands for P (logical processors), M (machine threads), and G (goroutines). The key point is that each logical processor (P) can only have one machine thread (M) running on it at any time. And for a goroutine (G) to run, it needs to be attached to a thread (M).

![alt text](https://victoriametrics.com/blog/go-sync-pool/sync-pool-pmg-model.webp)

- If you’ve got n logical processors (P), you can run up to n goroutines in parallel, as long as you’ve got at least n machine threads (M) available.

- At any one time, only one ``goroutine (G) ``can run on a ``single processor (P)``. So, when a P1 is busy with a G, no other G can run on that P1 until the current G either gets blocked, finishes up, or something else happens to free it up.

- But the thing is, a ``sync.Pool`` in Go isn’t just one big pool, it’s actually made up of several ’local’ pools, with each one tied to a specific processor context, or P, that Go’s runtime is managing at any given time.

![alt text](https://victoriametrics.com/blog/go-sync-pool/sync-pool-locals.webp)



## Best practices

- When instantiating ``sync.Pool`` give it a ``New`` member  variable that is ``thread-safe`` when called.

- When we receive an instanc from ```Get`` make no assumptions regarding the state of the object we receive back.

- Make sure to call ``Put`` when finished with the object, Otherwise the ``Pool`` is useless. Usually done with ``defer``.

- Objects in the pool must be roughly uniform in makeup.

- Best to use for objects that have ``rapid dispose`` after initializing and when the construction of these objects ``can negaitvely impact memory``.

- If code utilizes the ``Pool`` for requiring things that are not roughly homogenous, we may spend more time converting what we retrieved from the ``Pool``. 

- For slices of random and variable lenght , Pool is not going to help much.


## Further reading

[go-sync-pool](https://victoriametrics.com/blog/go-sync-pool/index.html)