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

