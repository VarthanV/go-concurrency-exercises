# Concurrency at scale

## Error Propogation

- In concurrent code and especially distributed systems its both easy for something to go wrong in our system and difficult to understand.

- Error indicate that our system has entered a state in which it cannot fulfill an operation that a user explicitly or implicity requested. Because of this , it need to relay a few pieces of critical information

**What happened?**:
- This is the part of the error that contains information about what happened. Eg, disk full ,socket closed or creential expired. This information is likely to be generated implicitly by whatever is generated the errors,although you can probably decorate with some context that will help the user.

**When and where it occured?**

- Errors should always contain a complete stack trace starting with how the call was initated and ending with where the error was instantiated. The stack trace should not be contained in the error message but should be easily accesible when handking the error up the stack.

- The error should contain information regarding the context is running within. For example in a distributed system it should have an way of identfying in what machine occured on.

- In addition, the error should contain the time on the machine the error was instantiated on , in UTC.


**A friendly user-facing message**

- The message that gets displayed should be customized to suit your system and its users.

- A friendly message is a human-centric gives some indication of whether the issue is transitory and should be about one line of text.

**How the user can get more information**

- At some point we will likely want to know in detail what happened when the error occured.

- Errors that are presented to users should provide an ID that can be cross referenced to the corresponding log that displays the full information of the error- time the error occured(not the time it was logged), the stacktrace everything we stuffed when it was created, it can also be helpful to include a ``hash`` of stack trace to aid in aggregating like issues in bug trackers.

- It is possible to place all error into one of two categories
    - Bugs
    - Known edge cases (eg. broken network connections, failed disk writes etc)

- Imagine a larger system with multiple modules

> CLI Component -> Intermediary Component -> Low Level Component

- Lets say an error occurs in the low level component, and we have created a well formed error, the error might be well formed in the context of the low level component but when pushing upwards the stack , it might not make sense to other components.

- It is crucial to wrap errors to have a clear cut understanding of where things went wrong for faster debugging.

- All the errors should be logged with as much as information as it is available.

- When our user facing code receives a well-formed error, we can be confident that at all levels in our code , care was taken to craft the message and we can simply log it out and print the messaage.

- When we get a malformed error , it is crucial to show a user friendly message,Let's say an error occurs internally with the message
> unexpected EOF

- The user might not understand this message.

- We have to attach a ``logID`` to the user for further resolvation of this issue.

- Typical well structured error

```go
type MyError struct {
    Inner error
    Message string
    StackTrace string
    Misc map[string]interface{}
}
```

