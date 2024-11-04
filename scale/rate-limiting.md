# Rate Limiting

- Constraints the number of times some kind of resource is accessed to some finite number of unit per time.

Eg: This API can be hit only 10times in one minute window

- The resource can be anything; API conncection , disk reads/writes, network packets,errors.

- By rate limiting a system we prevent entire classes of attack vector against the system.

- If not rate limited , Our log file might fill fastly and disk space will be filled causing the entire app to crash and no one can use it.

- Most rate limiting is done by utilizing an algorithm called the ``token buket``.

##  Token bucket

- Let's assume to utilize a resource we have an ``access token`` for the resource, Without the token the request is denied.

- Now imagine these tokens are stored in a bucket waiting to be retrieved for usage. This bucket has a depth of ``d``, which indicates it can hold ``d`` access tokens at a time. For eg if the bucket has d=5 it can hold 5 tokens.

- Now everytime we need to access a resource we reach into the bucket and remove a token. if the bucket contains five tokens and we can access the resource 5 times.

- Burstiness means how many requests we can make when the bucket is full.
