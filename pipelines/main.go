package main

import "fmt"

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

func main() {
	basicPipeline()
	WebScrapperPipelineDriver()
}
