# Timeouts and Cancellation

- When working with concurrent code , timeouts and cancellation are going to turn up frequently.

- Timeouts are crucial to creating a system with behaviour you can understand.

- Cancellation is one natural response to timeout.

## Reasons we might want to support time-outs

## System Saturation
- If our system is saturated (i.e. if its ability to process requests is at its capacity) we may want requests at the edge of our systems to timeout rather than take a long time to field them.

- If the request is unlikely to be repeated when it is timed out.

- If we don't have the resources to store the requests (eg. memory for in-memory queues,disk space for persisted queues).

- If need for request or the data it is sending will go stale. If the request is likely to be repeated system will develop an overhead from accepting and timing out requests. This can lead to death spirtal if over head becomes greater than our systems capacity.

## Stale Data

- Sometimes data has an window before which it needs to be processed, before more relavant data is available or need to process the data has expired. If the concurrent process takes longer to process the data than window, we would want to timeout and cancel the request.

- If the window is know beforehand it would make sense to pass our concurrent process a ``context.Context`` created with ``context.WithDeadline`` and ``context.Timeout`` if the deadline is known else a parent must able to cancel it when needed

**Eg:** Lets say we are sending an OTP to user which is valid for 5mins it doesn't makes  sense to process the request after the deadline.

## Attempting to  Prevent deadlocks

- The timeout period is to not pinpoint a time frame for completion of process , In a system when the calls propogate down the line there is a possiblity of ``deadlock`` , 
- The timeout is to unlock as soon as possible when a livelock occurs, because when a deadlock happens it can be fixed only by restart.

- It is much less of a overhead to fix a livelock than a dealock.

## Handling cancellation gracefully

- We need to consolidate the reasons why a concurrent process might be cancelled

**Timeouts**
    A timeout is an implicit cancellation

**User Intervention**
    - For a good user experience it is advisable to start long running processes concurrently and send report status back to the user at a polling interval or allow the users to query for status as they fit it.

    - When there are user facing concurrent OP's it also sometime necessary for allowing user to cancel the process.

**Parent  Cancellation**
    - If any kind of parent of concurrent operation - human or otherwise stops a child of the parent will be cancelled.

**Replicated Requests**
    - We might wish to send data to multiple concurrent processes in an attempt to get a faster response from one of them , when the first one comesback we would want to cancel the rest of the process.

## Preemtability of a concurrent process

- Our concurrent must be as preemtible as possible when cancellation occurs.

- Consider below code and assume it is running its own goroutine

```go
var value interface{}

select {
    case <- done:
        return
    case value := valueStream
}

result := reallyLongCalculation(value)

select {
    case <- done:
    return 
    case resultStream <- result:
}

```
- ``reallyLongCalculation`` doesn't look to be preemtable , This means some signal attempts to cancel the goroutine while ``reallyLongCalculation`` is executing it could be a very long time before we acknowledge the cancellation and halt

- We must maximum make the long running process in a goroutine as preemtable.

- An easy way to do this is to breakup pieces of goroutine into smaller pieces.

- Should aim for all nonpreemtable atomic operations to complete in less time than perod deemed acceptable.

- When cancellation signal comes when doing operations like database insertion, modification of file etc we need to be able rollback fairly quickly.

## Duplicated messages

- Another issue to be concerned is with ``duplicated messages``.

- Let's say in pipeline we have 3 stage , A,B and C, There is a possiblity of Stage B to receive duplicate messages if cancellation message comes in after stage A has already sent its result to B.

- There are few ways to avoid sending duplicate messages, The easiest is to make it vanishingly unlikely that a parent gorouine will send cancellation signal after a child goroutine has already reported a result. This requires **bidirectional communication** between stages.

**Accept either first or last result reported**: If the algorithm allows it or the concurrent process is idempotent , we can simply allow for the possiblity of duplicate messages in downstream processes and choose whether or not accept the first or last message

**Poll the goroutine for communication**: We can use bidirectional communication with parent to explicity request permission to send your message.

Note: 
> Introduce cancellations and timeouts in the beginning stage itself rather last

## Heartbeats

- Heartbeats are way for concurrent processes to signal life to outside parties.

- They allow us get insights from our system and they make the testing system deterministic which otherwise cannot be

- There are 2 different types of heartbeat
    - Occur on a time interval
    - Occur at beginning of a unit of work

- Heartbeat that occur on a time interval are useful for concurrent code that might be waiting for something else to happen for it to process a unit of work, because we dont know when the work might come in , goroutine might be sitting around for a while for something to happen. Heartbeat is a way to signal everything is fine
