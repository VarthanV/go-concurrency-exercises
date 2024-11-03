# Queueing

- Sometimes it is useful to begin accepting work for our pipeline eventhough the pipeline is not yet ready for more. The process is called ``queuing``

- This means is that once our stage has completed some work , it stores the result in a temporary location in memory so that the other stages can retrieve it later and our stage doesn't need to hold reference to it. This can be acheived via ``Buffered channels``.

- It is usually one the **last technique** we want to employ in **optimizing** our programs.

- Adding queue prematurely can hide synchornization issues such as **deadlocks and livelocks** and further we converge the program to correctness we need to find more or less queuing.

- Queuing will almost never speed up the total runtime of the program it will only the program to behave differently.

- Lets say we have a following pipeline
```go
    p := processRequest(done, acceptConnection(done,httpHandler))
```

- Here the pipeline doesn't exit until its cancelled and the stage that is accepting connections  doesnt stop accepting connection until the pipeline is cancelled.

- In this scenario we doesn't want our program to begin timing out because our ``processRequest`` stage was blocking our ``acceptConnection`` stage. We want our ``acceptConnection`` stage to be unblocked as much as possible. Otherwise the users of our program might being seen their requests almost denied together.

- The utility of queue isn't that the ``runtime of one of stages`` has been reduced , but rather the time it is in a ``blocking`` state is reduce. This allows the stage to continue doing its job.

- In this example , users would likely experience a lag in this requests but they wouldn't be denied service altogether.

- The true utility of the queues is to ``decouple stages`` so that the runtime of one stage has no impact on the runtime of other.

- Decouping stages in this manner then cascades or alter the runtime behavoir of the system as a whole, which can be either good or bad depending on your system.

## Tuning of Queue

-  The question naturally arises where the queues should be placed ? What should be the size of the buffer be?

- The answer depend on the nature of the pipeline.

- Situations where queue can increase overall performance of the system
    - If batching requests in a stage saves time.
    - If delays in a stage produce feedback loop into the system (In slow down in one stage cascades dropping in other stages).

## Batching requests in a stage saves time

- One example of the first situation is a stage that buffers input in something faster(eg.memory) than it   is designed to send to slower destination (eg.disk)

- This is ofcourse the purpose of Go's ``bufio`` package.

- The writes are queued internally into a buffer until sufficient chunk has been accumulated and then the chunk is written out. The process is called ``chunking``.

- Chunking is faster because ``bytes.Buffer`` must grow its allocated memory to accomdate the bytes it must store.

- For various reasons , growing memory is expensive therefore less times we have to grow the more efficient our system as a whole will perform. 

- Some examples of chunking
    - Opening db transactions
    - Calculating message checksums
    - Allocating contigous space


## Delays in a stage produce feedback loop into the system

- This stage where a delay in stage causes more input to the pipeline is a little more difficult to spot, but also more important because it can lead to a systematic collapse our upstream systems.

- The idea is often referred to as ``negative feedback loop``. This is because a recurrent relation exists between the pipeline and its upstream systems; The rate at which upstream stages or systems submit new request is linked to the efficieny of pipeline.

- If the efficiency of the pipeline drops below a certain threshold, the system upstream for the pipeline begin increasing their inputs into the pipeline which cause the pipeline to lose more effciency and death spiral begins. Without some sort of fail safe mechanism the pipeline system will never recover.

- By introducing a queue at the entrance of the pipeline we can break the feedback loop at the cost of creating lag for the requests.

- From the perspective of the caller into the pipeline , the request appears to be processing , but taking a very long time. As long as the caller doesn't time out our pipeline will remain stable . If caller timeout we need to be ready to support for check of readiness when dequing . We don't create ourselves feedback loop by processing dead requests.

- Queuing should be implemented either
    - At the entrance o the pipeline
    - Stages where batching will lead to higher efficiency.

