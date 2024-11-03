# Context

- As we have seen the idiom of ``done`` channel which flows through the program and cancels all blocking concurrent operations.

- It would be useful if we could communicate extra information alongside the simple notificiation to cancel; like why cancellation was occuring whether or not our fn has deadline by which it needs to complete.

- Context is a type that through functions like a done channel does.

- If we use ``context`` package each fn downstream from our tpo level concurrent call would take in ``Context`` as its first argument.

- Theres a ``Done`` method which returns a channel thats closed when our fn is to be preempted.

- A ``Deadline`` fn to indicate if a goroutine will be cancelled after a certain time.

- A ``Err`` method that will return non-nil if the goroutine was cancelled.

- The ``Value`` method to get value out of the context.

- The request specific information can be passed along the Context.

- The context package serves two primary  purposes
    - To provide an API for cancelling branches of the call graph.

    - To provide a data-bag for transporting request-scoped data through call graph.

- Cancellation has three aspects

    - A goroutine's parent may want to cancel it.

    - A goroutine may want to cancel its children

    - Any blocking operation with a goroutine need to preempted so that it may be cancelled.


- In the downstream , We cannot mutate the state of the underlying structure.

- The function that accepts the ``Context`` cannot cancel it. This protects fn up the call stack from children cancelling the context.

```go
func WithCancel(parent Context) (ctx Context , cancel CancelFunc)

func WithDeadline(parent context , deadline time.Time) (Context,CancelFunc)

func WithTimeout(parent Context , timeout time.Duration) (Context , CancelFunc)
```

- All these functions take in a Context and return new one. 

- ``WithCancel``: Returns a new ``Context`` that closes its done channel when the returned ``Cancel`` function is called.

- ``WithDeadline``: Returns a new ``Context`` that closes its done channel when the machine's clock advances past the given deadline.

- ``WithTimeout``: Returns a new ``Context`` that closes it done channel after the given timeout duration 

## WithTimeOut usage

- Sets a timeout duration from the moment it is called. Once the timeout expires, the context is canceled.

- Should use ``WithTimeout`` when  want an operation to be canceled if it takes longer than a specific amount of time.

- Useful for tasks where you know the maximum duration that should be allowed, like making API calls, database queries, or performing operations that should finish within a set timeframe.

## WithDeadline

- Sets a specific deadline (absolute time) for when the context should be canceled, regardless of when it was created.

- Use ``WithDeadline`` when have a specific end time in mind, such as waiting for the end of a business day or a scheduled event.

- Useful in cases where the exact time matters more than the duration.

- At the top of the asynchornous call-graph we probably wont have a been passed a ``Context``. The context package provides with two functions to create empty instances of the context

```go
func Background() Context
func TODO() Context
```

- ``Background`` simply returns a empty context.

- ``TODO`` is not meant for production purpose , but it remains empty context.

- When we are using either ``WithDeadline`` or ``WithTimeout`` we can get the deadline from the ``Deadline``  function call


```go
// A Context carries a deadline, cancellation signal, and request-scoped values
// across API boundaries. Its methods are safe for simultaneous use by multiple
// goroutines.
type Context interface {
    // Done returns a channel that is closed when this Context is canceled
    // or times out.
    Done() <-chan struct{}

    // Err indicates why this context was canceled, after the Done channel
    // is closed.
    Err() error

    // Deadline returns the time when this Context will be canceled, if any.
    Deadline() (deadline time.Time, ok bool)

    // Value returns the value associated with key or nil if none.
    Value(key interface{}) interface{}
}

```

- The qualifications of for a context key:

- The key we use must satisfy Go's notion of ``comparability`` that is , the equality operators ``==`` and ``!=`` need to return correct results when used.

- Values returned must be safe to access from multiple goroutines.

- It i s recommended to define a custom key-type in the package. As long as the other packages do the same, this prevents ``collisions`` within the Context.

- Although the context concept is well received by golang in the aspect of cancellation there are some minor things which we need to keep in mind when passing the data down the context.

1) The data should transit process or API boundaries
    If you generate the data in your process' memory its probably not a good candidate to be request scoped data unless we also pass it across a API boundary.

2) The data should be immutable
    If it is not , then we by definition what we are storing did not come from the request.

3) The data should trend toward simple types
    If request scoped data is meant to transit process and API boundaries , it much easier for other side to pull data out

4) The data should be data not types with methods
    Operations are logic and belong to things consuming this data

5) The data should help decorate options , not drive thme
    The algorithm should behave consistently irrespectively of what is passed in context.


| Data  | 1  | 2  | 3  | 4  | 5  |
|---|---|---|---|---|---|
|  Request ID | ✅ |  ✅ | ✅  | ✅  |  ✅ |
|  UserID | ✅  |  ✅ | ✅  | ✅  |   |
|  URL | ✅  | ✅  |   |   |   |
|  API server Connection |   |   |   |   |   |
|  Auth Token |✅ | ✅  |   ✅|  ✅ |   |
| Request Token | ✅| ✅| ✅| | |





