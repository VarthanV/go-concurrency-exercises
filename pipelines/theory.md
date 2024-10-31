# Pipelines

- Pipeline is a tool which we can use as form of abstraction in our system.

- It is very powerful tool when the program needs to process streams or batches of data.

- A pipeline is nothing more than a serious of things that take data in , perform operation on it and pass the data back out.

- We call each of these operations a ``stage`` of the ``pipeline``.

- By using a pipeline we seperate the concerns of each stage which provides numerous benefits.

- We can modify stages independent of one another , can mix and match how stages are combined independent of the modifying stages.

```go
func basicPipeline() {
	mutliply := func(values []int, multiplier int) []int {
		result := make([]int, 0, len(values))
		for _, val := range values {
			result = append(result, val*multiplier)
		}
		return result
	}

	add := func(values []int, additive int) []int {
		result := make([]int, 0, len(values))
		for _, val := range values {
			result = append(result, val+additive)
		}
		return result
	}

	// Driver code
	val := []int{1, 2, 3, 4}

	for _, v := range mutliply(add(val, 2), 1) {
		fmt.Println(v)
	}
}
```

## Properties of a Pipeline Stage

- A stage consumes and returns the same type

- A stage must be refied (Treat functions as first class citizens) by the language so that it may be passed around.

**Batch Processing**: Take a slice of data and returning a slice of data

**Stream Processing**: They operate on chunks of data all at once instead of one discrete value at a time.This means stage receives and emits one element at  a time.
