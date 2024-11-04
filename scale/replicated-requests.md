# Replicated Requests

- For some applications receiving a response as quickly as possible is the top priority.

- For example , the app is servicing a users HTTP request or retrieving a replicated blob of data.

- We can make a trade-off in these instances, we can replicate the request to multiple handlers  and when the first response comes we can immediately return the result. The downside is we have to utilize resources to keep multiple copies of handlers running

- If the replication is done in memory it might not be costly,but if the replication handler requires replicating process , servers or even data centers this can become quite costly.

- Replicate requests to handlers that have ``different runtime conditions``.

Eg: Lets say we are writing an weather app and we have two different providers to give same weather data we can replicate the request to these providers and return the first received response